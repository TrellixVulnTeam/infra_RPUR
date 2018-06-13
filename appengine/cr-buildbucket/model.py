# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import random

from components import auth
from components import datastore_utils
from google.appengine.ext import ndb

from google.appengine.ext.ndb import msgprop
from protorpc import messages

from proto import build_pb2
import buildtags

BEGINING_OF_THE_WORLD = datetime.datetime(2010, 1, 1, 0, 0, 0, 0)
BUILD_TIMEOUT = datetime.timedelta(days=2)

# If builds weren't scheduled for this duration on a given builder, the
# Builder entity is deleted.
BUILDER_EXPIRATION_DURATION = datetime.timedelta(weeks=4)

# Key in Build.parameters that specifies the builder name.
BUILDER_PARAMETER = 'builder_name'


class BuildStatus(messages.Enum):
  # A build is created, can be leased by someone and started.
  SCHEDULED = 1
  # Someone has leased the build and marked it as started.
  STARTED = 2
  # A build is completed. See BuildResult for more details.
  COMPLETED = 3


class BuildResult(messages.Enum):
  # A build has completed successfully.
  SUCCESS = 1
  # A build has completed unsuccessfully.
  FAILURE = 2
  # A build was canceled.
  CANCELED = 3


class FailureReason(messages.Enum):
  # Build failed
  BUILD_FAILURE = 1
  # Something happened within buildbucket.
  BUILDBUCKET_FAILURE = 2
  # Something happened with build infrastructure, but not buildbucket.
  INFRA_FAILURE = 3
  # A build-system rejected a build because its definition is invalid.
  INVALID_BUILD_DEFINITION = 4


class CancelationReason(messages.Enum):
  # A build was canceled explicitly, probably by an API call.
  CANCELED_EXPLICITLY = 1
  # A build was canceled by buildbucket due to timeout.
  TIMEOUT = 2


class CanaryPreference(messages.Enum):
  # The build system will decide whether to use canary or not
  AUTO = 1
  # Use the production build infrastructure
  PROD = 2
  # Use the canary build infrastructure
  CANARY = 3


class PubSubCallback(ndb.Model):
  """Parameters for a callack push task."""
  topic = ndb.StringProperty(required=True, indexed=False)
  auth_token = ndb.StringProperty(indexed=False)
  user_data = ndb.StringProperty(indexed=False)


class BucketState(ndb.Model):
  """Persistent state of a single bucket."""
  # If True, no new bulids may be leased for this bucket.
  is_paused = ndb.BooleanProperty()


class Build(ndb.Model):
  """Describes a build.

  Build key:
    Build keys are autogenerated, monotonically decreasing integers.
    That is, when sorted by key, new builds are first.
    Build has no parent.

    Build id is a 64 bits integer represented as a string to the user.
    - 1 highest order bit is set to 0 to keep value positive.
    - 43 bits are 43 lower bits of bitwise-inverted time since
      BEGINING_OF_THE_WORLD at 1ms resolution.
      It is good for 2**43 / 365.3 / 24 / 60 / 60 / 1000 = 278 years
      or 2010 + 278 = year 2288.
    - 16 bits are set to a random value. Assuming an instance is internally
      consistent with itself, it can ensure to not reuse the same 16 bits in two
      consecutive requests and/or throttle itself to one request per
      millisecond. Using random value reduces to 2**-15 the probability of
      collision on exact same timestamp at 1ms resolution, so a maximum
      theoretical rate of 65536000 requests/sec but an effective rate in the
      range of ~64k qps without much transaction conflicts. We should be fine.
    - 4 bits are 0. This is to represent the 'version' of the entity
      schema.

    The idea is taken from Swarming TaskRequest entity:
    https://code.google.com/p/swarming/source/browse/appengine/swarming/server/task_request.py#329
  """

  status = msgprop.EnumProperty(BuildStatus, default=BuildStatus.SCHEDULED)

  # A proto.common_pb2.Status corresponding to self.status.
  # This is needed to index builds by V2 status because status_v2->status_v1
  # function also depends on infra_failure_reason, i.e. without this it is
  # impossible to take a V2 status and translate it to Datastore query over
  # V1 status.
  status_v2 = ndb.ComputedProperty(lambda self: self._compute_v2_status())

  incomplete = ndb.ComputedProperty(
      lambda self: self.status != BuildStatus.COMPLETED)
  status_changed_time = ndb.DateTimeProperty(auto_now_add=True)
  update_time = ndb.DateTimeProperty(auto_now=True)

  # Creation time attributes.

  create_time = ndb.DateTimeProperty(auto_now_add=True)
  created_by = auth.IdentityProperty()
  # a generic way to distinguish builds.
  # Different buckets have different permissions.
  bucket = ndb.StringProperty(required=True)
  # property containing the ID of the LUCI project to which this build
  # belongs. Required for new builds, but older builds may not have it.
  project = ndb.StringProperty()
  # a list of tags, where each tag is a string
  # with ":" symbol. The first occurrence of ":" splits tag name and tag
  # value. Contains only tags specified by the build request. Old Build
  # entities do not have this field.
  initial_tags = ndb.StringProperty(repeated=True, indexed=False)
  # superset of initial_tags. May contain auto-added tags.
  tags = ndb.StringProperty(repeated=True)
  # immutable arbitrary build parameters.
  parameters = datastore_utils.DeterministicJsonProperty(json_type=dict)
  # PubSub message parameters for build status change notifications.
  pubsub_callback = ndb.StructuredProperty(PubSubCallback, indexed=False)
  # id of the original build that this build was derived from.
  retry_of = ndb.IntegerProperty()
  # Specifies whether canary of build infrastructure should be used for this
  # build.
  canary_preference = msgprop.EnumProperty(CanaryPreference, indexed=False)
  # If True, the build won't affect monitoring and won't be surfaced in
  # search results unless explicitly requested.
  experimental = ndb.BooleanProperty()

  # Lease-time attributes.

  # current lease expiration date.
  # The moment the build is leased, |lease_expiration_date| is set to
  # (utcnow + lease_duration).
  lease_expiration_date = ndb.DateTimeProperty()
  # None if build is not leased, otherwise a random value.
  # Changes every time a build is leased. Can be used to verify that a client
  # is the leaseholder.
  lease_key = ndb.IntegerProperty(indexed=False)
  # True if the build is currently leased. Otherwise False
  is_leased = ndb.ComputedProperty(lambda self: self.lease_key is not None)
  leasee = auth.IdentityProperty()
  never_leased = ndb.BooleanProperty()

  # Start time attributes.

  # a URL to a build-system-specific build, viewable by a human.
  url = ndb.StringProperty(indexed=False)
  # when the build started. Unknown for old builds.
  start_time = ndb.DateTimeProperty()
  # True if canary build infrastructure is used to run this build.
  # It may be None only in SCHEDULED state. Otherwise it must be True or False.
  # If canary_preference is CANARY, this field value does not have to be True,
  # e.g. if the build infrastructure does not have a canary.
  canary = ndb.BooleanProperty()

  # Completion time attributes.

  complete_time = ndb.DateTimeProperty()
  result = msgprop.EnumProperty(BuildResult)
  result_details = datastore_utils.DeterministicJsonProperty(json_type=dict)
  cancelation_reason = msgprop.EnumProperty(CancelationReason)
  failure_reason = msgprop.EnumProperty(FailureReason)

  # Swarming integration

  swarming_hostname = ndb.StringProperty()
  swarming_task_id = ndb.StringProperty()
  service_account = ndb.StringProperty()

  # LogDog integration

  logdog_hostname = ndb.StringProperty()
  logdog_project = ndb.StringProperty()
  logdog_prefix = ndb.StringProperty()

  def _pre_put_hook(self):
    """Checks Build invariants before putting."""
    super(Build, self)._pre_put_hook()
    is_started = self.status == BuildStatus.STARTED
    is_completed = self.status == BuildStatus.COMPLETED
    is_canceled = self.result == BuildResult.CANCELED
    is_failure = self.result == BuildResult.FAILURE
    is_leased = self.lease_key is not None
    assert (self.result is not None) == is_completed
    assert (self.cancelation_reason is not None) == is_canceled
    assert (self.failure_reason is not None) == is_failure
    assert not (is_completed and is_leased)
    assert (self.lease_expiration_date is not None) == is_leased
    assert (self.leasee is not None) == is_leased
    # no cover due to a bug in coverage (https://stackoverflow.com/a/35325514)

    tag_delm = buildtags.DELIMITER
    assert (not self.tags or
            all(tag_delm in t for t in self.tags))  # pragma: no cover
    assert self.create_time
    assert (self.complete_time is not None) == is_completed
    assert not is_started or self.start_time
    assert not self.start_time or self.start_time >= self.create_time
    assert not self.complete_time or self.complete_time >= self.create_time
    assert (not self.complete_time or not self.start_time or
            self.complete_time >= self.start_time)

    self.experimental = bool(self.experimental)
    self.initial_tags = sorted(set(self.initial_tags))
    self.tags = sorted(set(self.tags))

  def _compute_v2_status(self):
    build_v2 = build_pb2.Build()
    status_to_v2(self, build_v2)
    return build_v2.status

  def regenerate_lease_key(self):
    """Changes lease key to a different random int."""
    while True:
      new_key = random.randint(0, 1 << 31)
      if new_key != self.lease_key:  # pragma: no branch
        self.lease_key = new_key
        break

  def clear_lease(self):  # pragma: no cover
    """Clears build's lease attributes."""
    self.lease_key = None
    self.lease_expiration_date = None
    self.leasee = None


class BuildAnnotations(ndb.Model):
  """Stores annotation_pb2.Step of a build, if available.

  Available only for Swarmbucket builds, if we were able to retrieve it.
  Created on Build completion.

  Entity key:
    Parent is Build entity key.
    ID is 1.
  """

  ENTITY_ID = 1

  # root annotation_pb2.Step in binary format.
  annotation_binary = ndb.BlobProperty(compressed=True)
  # where the annotations_binary came from.
  annotation_url = ndb.StringProperty(indexed=False)

  @classmethod
  def key_for(cls, build_key):  # pragma: no cover
    return ndb.Key(cls, cls.ENTITY_ID, parent=build_key)


class Builder(ndb.Model):
  """A builder in a bucket.

  Used internally for metrics.
  Registered automatically by scheduling a build.
  Unregistered automatically by not scheduling builds for
  BUILDER_EXPIRATION_DURATION.

  Entity key:
    No parent. ID is a string with format "{project}:{bucket}:{builder}".
  """

  # Last time we received a valid build scheduling request for this builder.
  # Probabilistically updated by services.py, see its _should_update_builder.
  last_scheduled = ndb.DateTimeProperty()


class TagIndexEntry(ndb.Model):
  """A single entry in a TagIndex, references a build."""
  created_time = ndb.DateTimeProperty(auto_now_add=True)
  # ID of the build.
  build_id = ndb.IntegerProperty(indexed=False)
  # Bucket of the build.
  bucket = ndb.StringProperty(indexed=False)


class TagIndex(ndb.Model):
  """A custom index of builds by a tag.

  Entity key:
    Entity id is a build tag in the same "<key>:<value>" format that builds use.
    TagIndex has no parent.
  """

  MAX_ENTRY_COUNT = 1000

  # if incomplete, this TagIndex should not be used in search.
  # It is set to True if there are more than MAX_ENTRY_COUNT builds
  # for this tag.
  permanently_incomplete = ndb.BooleanProperty()

  # entries is a superset of all builds that have the tag equal to the id of
  # this entity. It may contain references to non-existent builds or builds that
  # do not actually have this tag; such builds must be ignored.
  #
  # It is sorted by build id in descending order.
  entries = ndb.LocalStructuredProperty(
      TagIndexEntry, repeated=True, indexed=False)


_TIME_RESOLUTION = datetime.timedelta(milliseconds=1)
_BUILD_ID_SUFFIX_LEN = 20


def _id_time_segment(dtime):
  assert dtime
  assert dtime >= BEGINING_OF_THE_WORLD
  delta = dtime - BEGINING_OF_THE_WORLD
  now = int(delta.total_seconds() * 1000.)
  return (~now & ((1 << 43) - 1)) << 20


def create_build_ids(dtime, count):
  """Returns a range of valid build ids, as integers and based on a datetime.

  See model.Build's docstring, "Build key" section.
  """
  # Build ID bits: "0N{43}R{16}V{4}"
  # where N is now bits, R is random bits and V is version bits.
  build_id = int(_id_time_segment(dtime) | (random.getrandbits(16) << 4))
  return [build_id - i * (1 << 4) for i in xrange(count)]


def build_id_range(create_time_low, create_time_high):
  """Converts a creation time range to build id range.

  Low/high bounds are inclusive/exclusive respectively, for both time and id
  ranges.
  """
  id_low = None
  id_high = None
  if create_time_low is not None:  # pragma: no branch
    # convert inclusive to exclusive
    id_high = _id_time_segment(create_time_low - _TIME_RESOLUTION)
  if create_time_high is not None:  # pragma: no branch
    # convert exclusive to inclusive
    id_low = _id_time_segment(create_time_high - _TIME_RESOLUTION)
  return id_low, id_high


status_to_v2 = None


def set_status_to_v2(fn):
  global status_to_v2
  assert status_to_v2 is None
  status_to_v2 = fn
