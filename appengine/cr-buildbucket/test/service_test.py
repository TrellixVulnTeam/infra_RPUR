# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import contextlib
import datetime
import json

from components import auth
from components import utils
from google.appengine.ext import ndb
from testing_utils import testing
import mock

from proto import build_pb2
from proto.config import service_config_pb2
from test import test_util
from test.test_util import future
import config
import errors
import model
import notifications
import service
import user
import v2


class BuildBucketServiceTest(testing.AppengineTestCase):

  def __init__(self, *args, **kwargs):
    super(BuildBucketServiceTest, self).__init__(*args, **kwargs)
    self.test_build = None

  def setUp(self):
    super(BuildBucketServiceTest, self).setUp()
    user.clear_request_cache()

    self.current_identity = auth.Identity('service', 'unittest')
    self.patch(
        'components.auth.get_current_identity',
        side_effect=lambda: self.current_identity
    )
    self.patch('user.can_async', return_value=future(True))
    self.now = datetime.datetime(2015, 1, 1)
    self.patch('components.utils.utcnow', side_effect=lambda: self.now)

    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg(
            '''
            name: "luci.chromium.try"
            acls {
              role: READER
              identity: "anonymous:anonymous"
            }
            swarming {
              hostname: "chromium-swarm.appspot.com"
              builders {
                name: "linux"
                build_numbers: YES
                recipe {
                  repository: "https://example.com"
                  name: "recipe"
                }
              }
            }
            '''
        ),
    )
    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg(
            '''
            name: "master.chromium"
            acls {
              role: READER
              identity: "anonymous:anonymous"
            }
            '''
        ),
    )

    self.patch('swarming.cancel_task_async', return_value=future(None))

    self.test_build = mkBuild()
    self.patch(
        'google.appengine.api.app_identity.get_default_version_hostname',
        autospec=True,
        return_value='buildbucket.example.com'
    )

    self.patch(
        'notifications.enqueue_tasks_async',
        autospec=True,
        return_value=future(None)
    )
    self.patch(
        'bq.enqueue_pull_task_async', autospec=True, return_value=future(None)
    )
    self.patch(
        'config.get_settings_async',
        autospec=True,
        return_value=future(service_config_pb2.SettingsCfg())
    )

    self.patch('search.TagIndex.random_shard_index', return_value=0)

  def mock_cannot(self, action, bucket_id=None):

    def can_async(requested_bucket_id, requested_action, _identity=None):
      match = (
          requested_action == action and
          (bucket_id is None or requested_bucket_id == bucket_id)
      )
      return future(not match)

    # user.can_async is patched in setUp()
    user.can_async.side_effect = can_async

  def put_many_builds(self, count=100, tags=None, **kwargs):
    tags = tags or []
    builds = []
    for _ in xrange(count):
      builds.append(mkBuild(**kwargs))
      self.now += datetime.timedelta(seconds=1)
    ndb.put_multi(builds)
    return builds

  #################################### GET #####################################

  def test_get(self):
    self.test_build.put()
    build = service.get_async(self.test_build.key.id()).get_result()
    self.assertEqual(build, self.test_build)

  def test_get_nonexistent_build(self):
    self.assertIsNone(service.get_async(42).get_result())

  def test_get_with_auth_error(self):
    self.mock_cannot(user.Action.VIEW_BUILD)
    self.test_build.put()
    with self.assertRaises(auth.AuthorizationError):
      service.get_async(self.test_build.key.id()).get_result()

  ################################### CANCEL ###################################

  def test_cancel(self):
    self.test_build.put()
    build = service.cancel(self.test_build.key.id(), human_reason='nope')
    self.assertEqual(build.status, model.BuildStatus.COMPLETED)
    self.assertEqual(build.status_changed_time, utils.utcnow())
    self.assertEqual(build.complete_time, utils.utcnow())
    self.assertEqual(build.result, model.BuildResult.CANCELED)
    self.assertEqual(
        build.cancelation_reason, model.CancelationReason.CANCELED_EXPLICITLY
    )
    self.assertEqual(
        build.cancel_reason_v2,
        build_pb2.CancelReason(
            message='nope',
            canceled_by=self.current_identity.to_bytes(),
        ),
    )

  def test_cancel_is_idempotent(self):
    self.test_build.put()
    service.cancel(self.test_build.key.id())
    service.cancel(self.test_build.key.id())

  def test_cancel_started_build(self):
    self.lease()
    self.start()
    service.cancel(self.test_build.key.id())

  def test_cancel_nonexistent_build(self):
    with self.assertRaises(errors.BuildNotFoundError):
      service.cancel(1)

  def test_cancel_with_auth_error(self):
    self.test_build.put()
    self.mock_cannot(user.Action.CANCEL_BUILD)
    with self.assertRaises(auth.AuthorizationError):
      service.cancel(self.test_build.key.id())

  def test_cancel_completed_build(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.SUCCESS
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()
    with self.assertRaises(errors.BuildIsCompletedError):
      service.cancel(self.test_build.key.id())

  @mock.patch('swarming.cancel_task_transactionally_async', autospec=True)
  def test_cancel_swarmbucket_build(self, cancel_task_async):
    cancel_task_async.return_value = future(None)
    self.test_build.swarming_hostname = 'chromium-swarm.appspot.com'
    self.test_build.swarming_task_id = 'deadbeef'
    self.test_build.put()
    service.cancel(self.test_build.key.id())
    cancel_task_async.assert_called_with(
        'chromium-swarm.appspot.com', 'deadbeef'
    )

  def test_cancel_result_details(self):
    self.test_build.put()
    result_details = {'message': 'bye bye build'}
    build = service.cancel(
        self.test_build.key.id(), result_details=result_details
    )
    self.assertEqual(build.result_details, result_details)

  def test_peek(self):
    self.test_build.put()
    builds, _ = service.peek(bucket_ids=[self.test_build.bucket_id])
    self.assertEqual(builds, [self.test_build])

  def test_peek_multi(self):
    self.test_build.key = ndb.Key(model.Build, 10)
    self.test_build.put()
    # We test that peek returns builds in decreasing order of the build key. The
    # build key is derived from the inverted current time, so later builds get
    # smaller ids. Only exception: if the time is the same, randomness decides
    # the order. So artificially create an id here to avoid flakiness.
    build2 = mkBuild(id=self.test_build.key.id() - 1,)
    build2.put()
    builds, _ = service.peek(
        bucket_ids=[self.test_build.bucket_id, 'chromium/ci']
    )
    self.assertEqual(builds, [self.test_build, build2])

  def test_peek_with_paging(self):
    self.put_many_builds()
    first_page, next_cursor = service.peek(
        bucket_ids=[self.test_build.bucket_id], max_builds=10
    )
    self.assertTrue(first_page)
    self.assertTrue(next_cursor)

    second_page, _ = service.peek(
        bucket_ids=[self.test_build.bucket_id], start_cursor=next_cursor
    )

    self.assertTrue(all(b not in second_page for b in first_page))

  def test_peek_with_bad_cursor(self):
    self.put_many_builds()
    with self.assertRaises(errors.InvalidInputError):
      service.peek(bucket_ids=[self.test_build.bucket_id], start_cursor='abc')

  def test_peek_without_buckets(self):
    with self.assertRaises(errors.InvalidInputError):
      service.peek(bucket_ids=[])

  def test_peek_with_auth_error(self):
    self.mock_cannot(user.Action.SEARCH_BUILDS)
    self.test_build.put()
    with self.assertRaises(auth.AuthorizationError):
      service.peek(bucket_ids=[self.test_build.bucket_id])

  def test_peek_does_not_return_leased_builds(self):
    self.test_build.put()
    self.lease()
    builds, _ = service.peek([self.test_build.bucket_id])
    self.assertFalse(builds)

  #################################### LEASE ###################################

  def lease(self, lease_expiration_date=None):
    if not (self.test_build.key and self.test_build.key.get()):
      self.test_build.put()
    success, self.test_build = service.lease(
        self.test_build.key.id(),
        lease_expiration_date=lease_expiration_date,
    )
    return success

  def test_lease(self):
    expiration_date = utils.utcnow() + datetime.timedelta(minutes=1)
    self.assertTrue(self.lease(lease_expiration_date=expiration_date))
    self.assertTrue(self.test_build.is_leased)
    self.assertGreater(self.test_build.lease_expiration_date, utils.utcnow())
    self.assertEqual(self.test_build.leasee, self.current_identity)

  def test_lease_build_with_auth_error(self):
    self.mock_cannot(user.Action.LEASE_BUILD)
    build = self.test_build
    build.put()
    with self.assertRaises(auth.AuthorizationError):
      self.lease()

  def test_cannot_lease_a_leased_build(self):
    build = self.test_build
    build.put()
    self.assertTrue(self.lease())
    self.assertFalse(self.lease())

  def test_cannot_lease_a_nonexistent_build(self):
    with self.assertRaises(errors.BuildNotFoundError):
      service.lease(build_id=42)

  def test_leasing_regenerates_lease_key(self):
    orig_lease_key = 42
    self.lease()
    self.assertNotEqual(self.test_build.lease_key, orig_lease_key)

  def test_cannot_lease_completed_build(self):
    build = self.test_build
    build.status = model.BuildStatus.COMPLETED
    build.result = model.BuildResult.SUCCESS
    build.complete_time = utils.utcnow()
    build.put()
    self.assertFalse(self.lease())

  ################################### UNELASE ##################################

  def test_reset(self):
    self.lease()
    build = service.reset(self.test_build.key.id())
    self.assertEqual(build.status, model.BuildStatus.SCHEDULED)
    self.assertEqual(build.status_changed_time, utils.utcnow())
    self.assertIsNone(build.lease_key)
    self.assertIsNone(build.lease_expiration_date)
    self.assertIsNone(build.leasee)
    self.assertIsNone(build.canary)
    self.assertTrue(self.lease())

  def test_reset_is_idempotent(self):
    self.lease()
    build_id = self.test_build.key.id()
    service.reset(build_id)
    service.reset(build_id)

  def test_reset_completed_build(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.SUCCESS
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()

    with self.assertRaises(errors.BuildIsCompletedError):
      service.reset(self.test_build.key.id())

  def test_cannot_reset_nonexistent_build(self):
    with self.assertRaises(errors.BuildNotFoundError):
      service.reset(123)

  def test_reset_with_auth_error(self):
    self.lease()
    self.mock_cannot(user.Action.RESET_BUILD)
    with self.assertRaises(auth.AuthorizationError):
      service.reset(self.test_build.key.id())

  #################################### START ###################################

  def test_validate_malformed_url(self):
    with self.assertRaises(errors.InvalidInputError):
      service.validate_url('svn://sdfsf')

  def test_validate_relative_url(self):
    with self.assertRaises(errors.InvalidInputError):
      service.validate_url('sdfsf')

  def test_validate_nonstring_url(self):
    with self.assertRaises(errors.InvalidInputError):
      service.validate_url(123)

  def start(self, url=None, lease_key=None, canary=False):
    self.test_build = service.start(
        self.test_build.key.id(), lease_key or self.test_build.lease_key, url,
        canary
    )

  def test_start(self):
    self.lease()
    self.start(url='http://localhost', canary=True)
    self.assertEqual(self.test_build.status, model.BuildStatus.STARTED)
    self.assertEqual(self.test_build.url, 'http://localhost')
    self.assertEqual(self.test_build.start_time, self.now)
    self.assertTrue(self.test_build.canary)

  def test_start_started_build(self):
    self.lease()
    build_id = self.test_build.key.id()
    lease_key = self.test_build.lease_key
    url = 'http://localhost/'

    service.start(build_id, lease_key, url, False)
    service.start(build_id, lease_key, url, False)
    service.start(build_id, lease_key, url + '1', False)

  def test_start_non_leased_build(self):
    self.test_build.put()
    with self.assertRaises(errors.LeaseExpiredError):
      service.start(self.test_build.key.id(), 42, None, False)

  def test_start_completed_build(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.SUCCESS
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()
    with self.assertRaises(errors.BuildIsCompletedError):
      service.start(self.test_build.key.id(), 42, None, False)

  def test_start_without_lease_key(self):
    with self.assertRaises(errors.InvalidInputError):
      service.start(1, None, None, False)

  @contextlib.contextmanager
  def callback_test(self):
    self.test_build.key = ndb.Key(model.Build, 1)
    self.test_build.pubsub_callback = model.PubSubCallback(
        topic='projects/example/topics/buildbucket',
        user_data='hello',
        auth_token='secret',
    )
    self.test_build.put()
    yield
    notifications.enqueue_tasks_async.assert_called_with(
        'backend-default', [
            {
                'url':
                    '/internal/task/buildbucket/notify/1',
                'payload':
                    json.dumps({
                        'id': 1,
                        'mode': 'global',
                    }, sort_keys=True),
                'age_limit_sec':
                    model.BUILD_TIMEOUT.total_seconds(),
            },
            {
                'url':
                    '/internal/task/buildbucket/notify/1',
                'payload':
                    json.dumps({
                        'id': 1,
                        'mode': 'callback',
                    }, sort_keys=True),
                'age_limit_sec':
                    model.BUILD_TIMEOUT.total_seconds(),
            },
        ]
    )

  def test_start_creates_notification_task(self):
    self.lease()
    with self.callback_test():
      self.start()

  ################################## HEARTBEAT #################################

  def test_heartbeat(self):
    self.lease()
    new_expiration_date = utils.utcnow() + datetime.timedelta(minutes=1)
    build = service.heartbeat(
        self.test_build.key.id(),
        self.test_build.lease_key,
        lease_expiration_date=new_expiration_date
    )
    self.assertEqual(build.lease_expiration_date, new_expiration_date)

  def test_heartbeat_completed(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.CANCELED
    self.test_build.cancelation_reason = (
        model.CancelationReason.CANCELED_EXPLICITLY
    )
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()

    new_expiration_date = utils.utcnow() + datetime.timedelta(minutes=1)
    with self.assertRaises(errors.BuildIsCompletedError):
      service.heartbeat(
          self.test_build.key.id(),
          0,
          lease_expiration_date=new_expiration_date
      )

  def test_heartbeat_timed_out(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.CANCELED
    self.test_build.cancelation_reason = model.CancelationReason.TIMEOUT
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()

    new_expiration_date = utils.utcnow() + datetime.timedelta(minutes=1)
    exc_regex = (
        'Build was marked as timed out '
        'because it did not complete for 2 days'
    )
    with self.assertRaisesRegexp(errors.BuildIsCompletedError, exc_regex):
      service.heartbeat(
          self.test_build.key.id(),
          0,
          lease_expiration_date=new_expiration_date
      )

  def test_heartbeat_batch(self):
    self.lease()
    new_expiration_date = utils.utcnow() + datetime.timedelta(minutes=1)
    results = service.heartbeat_batch([
        {
            'build_id': self.test_build.key.id(),
            'lease_key': self.test_build.lease_key,
            'lease_expiration_date': new_expiration_date,
        },
        {
            'build_id': 42,
            'lease_key': 42,
            'lease_expiration_date': new_expiration_date,
        },
    ])

    self.assertEqual(len(results), 2)

    self.test_build = self.test_build.key.get()
    self.assertEqual(
        results[0], (self.test_build.key.id(), self.test_build, None)
    )

    self.assertIsNone(results[1][1])
    self.assertTrue(isinstance(results[1][2], errors.BuildNotFoundError))

  def test_heartbeat_without_expiration_date(self):
    self.lease()
    with self.assertRaises(errors.InvalidInputError):
      service.heartbeat(
          self.test_build.key.id(),
          self.test_build.lease_key,
          lease_expiration_date=None
      )

  ################################### COMPLETE #################################

  def succeed(self, **kwargs):
    self.test_build = service.succeed(
        self.test_build.key.id(), self.test_build.lease_key, **kwargs
    )

  def test_succeed(self):
    self.lease()
    self.start()
    self.succeed(result_details={'properties': {'foo': 'bar',}})
    self.assertEqual(self.test_build.status, model.BuildStatus.COMPLETED)
    self.assertEqual(self.test_build.status_changed_time, utils.utcnow())
    self.assertEqual(self.test_build.result, model.BuildResult.SUCCESS)
    self.assertIsNotNone(self.test_build.complete_time)

    out_props = model.BuildOutputProperties.key_for(self.test_build.key).get()
    self.assertEqual(
        test_util.msg_to_dict(out_props.properties), {'foo': 'bar'}
    )

  def test_succeed_timed_out_build(self):
    self.test_build.status = model.BuildStatus.COMPLETED
    self.test_build.result = model.BuildResult.CANCELED
    self.test_build.cancelation_reason = model.CancelationReason.TIMEOUT
    self.test_build.complete_time = utils.utcnow()
    self.test_build.put()
    with self.assertRaises(errors.BuildIsCompletedError):
      service.succeed(self.test_build.key.id(), 42)

  def test_succeed_is_idempotent(self):
    self.lease()
    self.start()
    build_id = self.test_build.key.id()
    lease_key = self.test_build.lease_key
    service.succeed(build_id, lease_key)
    service.succeed(build_id, lease_key)

  def test_succeed_with_new_tags(self):
    self.test_build.tags = ['a:1']
    self.test_build.put()
    self.lease()
    self.start()
    self.succeed(new_tags=['b:2'])
    self.assertEqual(self.test_build.tags, ['a:1', 'b:2'])

  def test_fail(self):
    self.lease()
    self.start()
    self.test_build = service.fail(
        self.test_build.key.id(), self.test_build.lease_key
    )
    self.assertEqual(self.test_build.status, model.BuildStatus.COMPLETED)
    self.assertEqual(self.test_build.status_changed_time, utils.utcnow())
    self.assertEqual(self.test_build.result, model.BuildResult.FAILURE)
    self.assertIsNotNone(self.test_build.complete_time)

  def test_fail_with_details(self):
    self.lease()
    self.start()
    result_details = {'transient_failure': True}
    self.test_build = service.fail(
        self.test_build.key.id(),
        self.test_build.lease_key,
        result_details=result_details
    )
    self.assertEqual(self.test_build.result_details, result_details)

  def test_complete_with_url(self):
    self.lease()
    self.start()
    url = 'http://localhost/1'
    self.succeed(url=url)
    self.assertEqual(self.test_build.url, url)

  def test_complete_not_started_build(self):
    self.lease()
    self.succeed()

  def test_completion_creates_notification_task(self):
    self.lease()
    self.start()
    with self.callback_test():
      self.succeed()

  ########################## RESET EXPIRED BUILDS ##############################

  def test_delete_many_scheduled_builds(self):
    self.test_build.put()
    completed_build = mkBuild(
        status=model.BuildStatus.COMPLETED,
        result=model.BuildResult.SUCCESS,
        complete_time=utils.utcnow() + datetime.timedelta(seconds=1),
    )
    completed_build.put()
    self.assertIsNotNone(self.test_build.key.get())
    self.assertIsNotNone(completed_build.key.get())
    service._task_delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.SCHEDULED
    )
    self.assertIsNone(self.test_build.key.get())
    self.assertIsNotNone(completed_build.key.get())

  def test_delete_many_started_builds(self):
    self.test_build.put()

    started_build = mkBuild(
        status=model.BuildStatus.STARTED,
        start_time=utils.utcnow(),
    )
    started_build.put()

    completed_build = mkBuild(
        status=model.BuildStatus.COMPLETED,
        result=model.BuildResult.SUCCESS,
        create_time=utils.utcnow(),
        complete_time=utils.utcnow(),
    )
    completed_build.put()

    service._task_delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.STARTED
    )
    self.assertIsNotNone(self.test_build.key.get())
    self.assertIsNone(started_build.key.get())
    self.assertIsNotNone(completed_build.key.get())

  def test_delete_many_builds_with_tags(self):
    self.test_build.tags = ['tag:1']
    self.test_build.put()

    service._task_delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.SCHEDULED, tags=['tag:0']
    )
    self.assertIsNotNone(self.test_build.key.get())

    service._task_delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.SCHEDULED, tags=['tag:1']
    )
    self.assertIsNone(self.test_build.key.get())

  def test_delete_many_builds_created_by(self):
    self.test_build.created_by = auth.Identity('user', 'nodir@google.com')
    self.test_build.put()
    other_build = mkBuild()
    other_build.put()

    service._task_delete_many_builds(
        self.test_build.bucket_id,
        model.BuildStatus.SCHEDULED,
        created_by='nodir@google.com'
    )
    self.assertIsNone(self.test_build.key.get())
    self.assertIsNotNone(other_build.key.get())

  def test_delete_many_builds_auth_error(self):
    self.mock_cannot(user.Action.DELETE_SCHEDULED_BUILDS)
    with self.assertRaises(auth.AuthorizationError):
      service.delete_many_builds(
          self.test_build.bucket_id, model.BuildStatus.SCHEDULED
      )

  def test_delete_many_builds_schedule_task(self):
    service.delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.SCHEDULED
    )

  def test_delete_many_completed_builds(self):
    with self.assertRaises(errors.InvalidInputError):
      service.delete_many_builds(
          self.test_build.bucket_id, model.BuildStatus.COMPLETED
      )

  @mock.patch('swarming.cancel_task_transactionally_async', autospec=True)
  def test_delete_many_swarmbucket_builds(self, cancel_task_async):
    cancel_task_async.return_value = future(None)
    self.test_build.swarming_hostname = 'swarming.example.com'
    self.test_build.swarming_task_id = 'deadbeef'
    self.test_build.put()

    service._task_delete_many_builds(
        self.test_build.bucket_id, model.BuildStatus.SCHEDULED
    )

    cancel_task_async.assert_called_with('swarming.example.com', 'deadbeef')

  ################################ PAUSE BUCKET ################################

  def test_pause_bucket(self):
    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg('name: "master.foo"'),
    )
    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg('name: "master.bar"'),
    )

    self.put_many_builds(5, bucket_id='chromium/master.foo')
    self.put_many_builds(5, bucket_id='chromium/master.bar')

    service.pause('chromium/master.foo', True)
    builds, _ = service.peek(['chromium/master.foo', 'chromium/master.bar'])
    self.assertEqual(len(builds), 5)
    self.assertTrue(all(b.bucket_id == 'chromium/master.bar' for b in builds))

  def test_pause_all_requested_buckets(self):
    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg('name: "master.foo"'),
    )
    self.put_many_builds(5, bucket_id='chromium/master.foo')

    service.pause('chromium/master.foo', True)
    builds, _ = service.peek(['chromium/master.foo'])
    self.assertEqual(len(builds), 0)

  def test_pause_then_unpause(self):
    bid = 'chromium/master.foo'
    self.test_build.bucket_id = bid
    self.test_build.put()

    config.put_bucket(
        'chromium',
        'a' * 40,
        test_util.parse_bucket_cfg('name: "master.foo"'),
    )

    service.pause(bid, True)
    service.pause(bid, True)  # Again, to cover equality case.
    builds, _ = service.peek([bid])
    self.assertEqual(len(builds), 0)

    service.pause(bid, False)
    builds, _ = service.peek([bid])
    self.assertEqual(len(builds), 1)

  def test_pause_bucket_auth_error(self):
    self.mock_cannot(user.Action.PAUSE_BUCKET)
    with self.assertRaises(auth.AuthorizationError):
      service.pause('chromium/no.such.bucket', True)

  def test_pause_invalid_bucket(self):
    config.get_bucket_async.return_value = future((None, None))
    with self.assertRaises(errors.InvalidInputError):
      service.pause('a/#', True)

  def test_pause_swarming_bucket(self):
    with self.assertRaises(errors.InvalidInputError):
      service.pause('chromium/try', True)

  ############################ UNREGISTER BUILDERS #############################

  def test_unregister_builders(self):
    model.Builder(
        id='chromium:try:linux_rel',
        last_scheduled=self.now - datetime.timedelta(weeks=8),
    ).put()
    service.unregister_builders()
    builders = model.Builder.query().fetch()
    self.assertFalse(builders)


def mkBuild(**kwargs):
  args = dict(
      id=model.create_build_ids(utils.utcnow(), 1)[0],
      proto=build_pb2.Build(),
      bucket_id='chromium/try',
      create_time=utils.utcnow(),
      created_by=auth.Identity('user', 'john@example.com'),
      canary_preference=model.CanaryPreference.PROD,
      parameters={
          model.BUILDER_PARAMETER: 'linux',
      },
      canary=False,
  )
  args.update(kwargs)
  return model.Build(**args)
