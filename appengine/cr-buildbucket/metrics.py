# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import logging

from google.appengine.ext import ndb

from components import utils
import gae_ts_mon

import buildtags
import model

# Override default target fields for app-global metrics.
GLOBAL_TARGET_FIELDS = {
    'job_name': '',  # module name
    'hostname': '',  # version
    'task_num': 0,  # instance ID
}

_SOURCE_TAG = 0  # read field the value from a tag
_SOURCE_ATTR = 1  # read field the value from an attr
# tuple item meanings are: value source, gae_ts_mon field type, value type.
_TAG_STR_FIELD = (_SOURCE_TAG, gae_ts_mon.StringField, str)
_ATTR_STR_FIELD = (_SOURCE_ATTR, gae_ts_mon.StringField, str)
_BUILD_FIELDS = {
    'bucket': _ATTR_STR_FIELD,
    'builder': _TAG_STR_FIELD,
    'canary': (_SOURCE_ATTR, gae_ts_mon.BooleanField, bool),
    'cancelation_reason': _ATTR_STR_FIELD,
    'failure_reason': _ATTR_STR_FIELD,
    'result': _ATTR_STR_FIELD,
    'status': _ATTR_STR_FIELD,
    'user_agent': _TAG_STR_FIELD,
}
_METRIC_PREFIX_PROD = 'buildbucket/builds/'
_METRIC_PREFIX_EXPERIMENTAL = 'buildbucket/builds-experimental/'

# Maximum number of concurrent counting/latency queries.
_CONCURRENT_QUERY_LIMIT = 100


def _default_field_value(name):
  return _BUILD_FIELDS[name][2]()


BUCKETER_24_HR = gae_ts_mon.GeometricBucketer(growth_factor=10**0.05)
BUCKETER_48_HR = gae_ts_mon.GeometricBucketer(growth_factor=10**0.053)
BUCKETER_5_SEC = gae_ts_mon.GeometricBucketer(growth_factor=10**0.0374)
BUCKETER_1K = gae_ts_mon.GeometricBucketer(growth_factor=10**0.031)


def _fields_for(build, field_names):
  """Returns field values for a build"""
  for f in field_names:
    if f not in _BUILD_FIELDS:
      raise ValueError('invalid field %r' % f)

  tags = None
  result = {}
  for f in field_names:
    src, _, typ = _BUILD_FIELDS[f]
    assert src in (_SOURCE_ATTR, _SOURCE_TAG)
    if src == _SOURCE_ATTR:
      val = getattr(build, f)
    else:
      if tags is None:
        tags = dict(buildtags.parse(t) for t in build.tags)
      val = tags.get(f)
    result[f] = typ(val or _default_field_value(f))
  return result


def _fields_for_fn(fields):
  assert all(f.name in _BUILD_FIELDS for f in fields)
  field_names = [f.name for f in fields]
  return lambda b: _fields_for(b, field_names)  # pragma: no cover


def _build_fields(*names):
  return [_BUILD_FIELDS[n][1](n) for n in names]


def _incrementer(metric_suffix, description, fields):
  """Returns a function that increments a counter metric.

  Metric fields must conform _BUILD_FIELDS.

  The returned function accepts a build.
  """

  def mk_metric(metric_prefix):
    return gae_ts_mon.CounterMetric(metric_prefix + metric_suffix, description,
                                    fields)

  metric_prod = mk_metric(_METRIC_PREFIX_PROD)
  metric_exp = mk_metric(_METRIC_PREFIX_EXPERIMENTAL)
  fields_for = _fields_for_fn(fields)

  def inc(build):  # pragma: no cover
    metric = metric_exp if build.experimental else metric_prod
    metric.increment(fields_for(build))

  return inc


def _adder(metric_suffix, description, fields, bucketer, units, value_fn):
  """Returns a function that adds a build value to a cumulative distribution.

  Metric fields must conform _BUILD_FIELDS.
  value_fn accepts a build.

  The returned function accepts a build.
  """

  def mk_metric(metric_prefix):
    return gae_ts_mon.CumulativeDistributionMetric(
        metric_prefix + metric_suffix,
        description,
        fields,
        bucketer=bucketer,
        units=units)

  metric_prod = mk_metric(_METRIC_PREFIX_PROD)
  metric_exp = mk_metric(_METRIC_PREFIX_EXPERIMENTAL)
  fields_for = _fields_for_fn(fields)

  def add(build):  # pragma: no cover
    metric = metric_exp if build.experimental else metric_prod
    metric.add(value_fn(build), fields_for(build))

  return add


def _duration_adder(metric_suffix, description, value_fn):
  return _adder(metric_suffix, description,
                _build_fields('bucket', 'builder', 'result', 'failure_reason',
                              'cancelation_reason', 'canary'), BUCKETER_48_HR,
                gae_ts_mon.MetricsDataUnits.SECONDS, value_fn)


inc_created_builds = _incrementer('created', 'Build creation',
                                  _build_fields('bucket', 'builder',
                                                'user_agent'))
inc_started_builds = _incrementer('started', 'Build start',
                                  _build_fields('bucket', 'builder', 'canary'))
inc_completed_builds = _incrementer(
    'completed',
    'Build completion, including success, failure and cancellation',
    _build_fields('bucket', 'builder', 'result', 'failure_reason',
                  'cancelation_reason', 'canary'))
inc_lease_expirations = _incrementer('lease_expired', 'Build lease expirations',
                                     _build_fields('bucket', 'builder',
                                                   'status'))
inc_leases = _incrementer('leases',
                          'Successful build leases or lease extensions',
                          _build_fields('bucket', 'builder'))

inc_heartbeat_failures = gae_ts_mon.CounterMetric(
    'buildbucket/builds/heartbeats', 'Failures to extend a build lease',
    []).increment

# requires the argument to have non-None create_time and complete_time.
add_build_cycle_duration = _duration_adder(  # pragma: no branch
    'cycle_durations', 'Duration between build creation and completion',
    lambda b: (b.complete_time - b.create_time).total_seconds())

# requires the argument to have non-None start_time and complete_time.
add_build_run_duration = _duration_adder(  # pragma: no branch
    'run_durations', 'Duration between build start and completion',
    lambda b: (b.complete_time - b.start_time).total_seconds())

# requires the argument to have non-None create_time and start_time.
add_build_scheduling_duration = _duration_adder(  # pragma: no branch
    'scheduling_durations', 'Duration between build creation and start',
    lambda b: (b.start_time - b.create_time).total_seconds())

BUILD_COUNT_PROD = gae_ts_mon.GaugeMetric(
    _METRIC_PREFIX_PROD + 'count', 'Number of pending/running prod builds',
    _build_fields('bucket', 'builder', 'status'))
BUILD_COUNT_EXPERIMENTAL = gae_ts_mon.GaugeMetric(
    _METRIC_PREFIX_EXPERIMENTAL + 'count',
    'Number of pending/running experimental builds',
    _build_fields('bucket', 'builder', 'status'))

MAX_AGE_SCHEDULED = gae_ts_mon.FloatMetric(
    'buildbucket/builds/max_age_scheduled',
    'Age of the oldest SCHEDULED build.',
    (_build_fields('bucket', 'builder') +
     [gae_ts_mon.BooleanField('must_be_never_leased')]),
    units=gae_ts_mon.MetricsDataUnits.SECONDS)
SEQUENCE_NUMBER_GEN_DURATION_MS = gae_ts_mon.CumulativeDistributionMetric(
    'buildbucket/sequence_number/gen_duration',
    'Duration of a sequence number generation in ms',
    [gae_ts_mon.StringField('sequence')],
    # Bucketer for 1ms..5s range
    bucketer=BUCKETER_5_SEC,
    units=gae_ts_mon.MetricsDataUnits.MILLISECONDS)
TAG_INDEX_INCONSISTENT_ENTRIES = gae_ts_mon.NonCumulativeDistributionMetric(
    'buildbucket/tag_index/inconsistent_entries',
    'Number of inconsistent entries encountered during build search',
    [gae_ts_mon.StringField('tag')],
    # We can't have more than 1000 entries in a tag index.
    bucketer=BUCKETER_1K)
TAG_INDEX_SEARCH_SKIPPED_BUILDS = gae_ts_mon.NonCumulativeDistributionMetric(
    'buildbucket/tag_index/skipped_builds',
    'Number of builds we fetched, but skipped',
    [gae_ts_mon.StringField('tag')],
    # We can't have more than 1000 entries in a tag index.
    bucketer=BUCKETER_1K)


@ndb.tasklet
def set_build_count_metric_async(bucket, builder, status, experimental):
  assert isinstance(bucket, basestring)
  assert isinstance(builder, basestring)
  q = model.Build.query(
      model.Build.bucket == bucket,
      model.Build.tags == 'builder:%s' % builder,
      model.Build.status == status,
      model.Build.experimental == experimental,
  )
  try:
    value = yield q.count_async()
  except Exception:  # pragma: no cover
    logging.exception('failed to count builds with query %s', q)
    return

  fields = {
      'bucket': bucket,
      'builder': builder,
      'status': str(status),
  }
  metric = BUILD_COUNT_EXPERIMENTAL if experimental else BUILD_COUNT_PROD
  metric.set(value, fields=fields, target_fields=GLOBAL_TARGET_FIELDS)


@ndb.tasklet
def set_build_latency(bucket, builder, must_be_never_leased):
  q = model.Build.query(
      model.Build.bucket == bucket,
      model.Build.tags == 'builder:%s' % builder,
      model.Build.status == model.BuildStatus.SCHEDULED,
      model.Build.experimental == False,
  )
  if must_be_never_leased:
    q = q.filter(model.Build.never_leased == True)
  else:
    # Reuse the index that has never_leased
    q = q.filter(model.Build.never_leased.IN((True, False)))

  q = q.order(model.Build.create_time)
  oldest_build = yield q.fetch_async(1, projection=[model.Build.create_time])

  if oldest_build:
    max_age = (utils.utcnow() - oldest_build[0].create_time).total_seconds()
  else:
    max_age = 0
  fields = {
      'bucket': bucket,
      'builder': builder,
      'must_be_never_leased': must_be_never_leased,
  }
  MAX_AGE_SCHEDULED.set(max_age, fields, target_fields=GLOBAL_TARGET_FIELDS)


# Metrics that are per-app rather than per-instance.
GLOBAL_METRICS = [
    MAX_AGE_SCHEDULED,
    BUILD_COUNT_PROD,
    BUILD_COUNT_EXPERIMENTAL,
]


def update_global_metrics():
  """Updates the metrics in GLOBAL_METRICS."""
  start = utils.utcnow()

  # Collect a list of counting/latency queries.
  count_query_queue = []
  latency_query_queue = []
  for key in model.Builder.query().iter(keys_only=True):
    _, bucket, builder = key.id().split(':', 2)
    latency_query_queue.extend([
        (bucket, builder, True),
        (bucket, builder, False),
    ])
    for status in (model.BuildStatus.SCHEDULED, model.BuildStatus.STARTED):
      for experimental in (False, True):
        count_query_queue.append((bucket, builder, status, experimental))

  # Process counting/latency queries with _CONCURRENT_QUERY_LIMIT workers.

  @ndb.tasklet
  def worker():
    while count_query_queue:
      item = count_query_queue.pop()
      yield set_build_count_metric_async(*item)
    while latency_query_queue:
      item = latency_query_queue.pop()
      yield set_build_latency(*item)

  for w in [worker() for _ in xrange(_CONCURRENT_QUERY_LIMIT)]:
    w.check_success()

  logging.info('global metric computation took %s', utils.utcnow() - start)
