# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime

from components import auth
from components import net
from components import utils
from google.appengine.ext import ndb
from testing_utils import testing
import mock

from proto import build_pb2
from proto import common_pb2
from proto import rpc_pb2
from proto.config import service_config_pb2
from test import test_util
from test.test_util import future, future_exception
import bbutil
import config
import creation
import errors
import model
import search
import swarming
import user
import v2


class CreationTest(testing.AppengineTestCase):
  test_build = None

  def setUp(self):
    super(CreationTest, self).setUp()
    user.clear_request_cache()

    self.current_identity = auth.Identity('service', 'unittest')
    self.patch(
        'components.auth.get_current_identity',
        side_effect=lambda: self.current_identity
    )
    self.patch('user.can_async', return_value=future(True))
    self.now = datetime.datetime(2015, 1, 1)
    self.patch('components.utils.utcnow', side_effect=lambda: self.now)

    self.chromium_try = test_util.parse_bucket_cfg(
        '''
        name: "luci.chromium.try"
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
    )
    config.put_bucket('chromium', 'a' * 40, self.chromium_try)
    self.patch('swarming.create_task_async', return_value=future(None))
    self.patch('swarming.cancel_task_async', return_value=future(None))

    self.test_build = model.Build(
        id=model.create_build_ids(self.now, 1)[0],
        bucket_id='chromium/try',
        create_time=self.now,
        parameters={
            model.BUILDER_PARAMETER:
                'linux',
            'changes': [{
                'author': 'nodir@google.com',
                'message': 'buildbucket: initial commit'
            }],
        },
        canary=False,
    )

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
        'config.get_settings_async',
        autospec=True,
        return_value=future(service_config_pb2.SettingsCfg())
    )

    self.patch('creation._should_update_builder', side_effect=lambda p: p > 0.5)

    self.patch('search.TagIndex.random_shard_index', return_value=0)

  def build_request(self, schedule_build_request_fields=None, **kwargs):
    schedule_build_request_fields = schedule_build_request_fields or {}
    sbr = rpc_pb2.ScheduleBuildRequest(**schedule_build_request_fields)
    sbr.builder.project = sbr.builder.project or 'chromium'
    sbr.builder.bucket = sbr.builder.bucket or 'try'
    sbr.builder.builder = sbr.builder.builder or 'linux'
    return creation.BuildRequest(schedule_build_request=sbr, **kwargs)

  def add(self, *args, **kwargs):
    br = self.build_request(*args, **kwargs)
    return creation.add_async(br).get_result()

  def test_add(self):
    builder_id = build_pb2.BuilderID(
        project='chromium',
        bucket='try',
        builder='linux',
    )
    build = self.add(dict(builder=builder_id))
    self.assertIsNotNone(build.key)
    self.assertIsNotNone(build.key.id())

    build = build.key.get()
    self.assertEqual(build.proto.builder, builder_id)
    self.assertEqual(
        build.proto.created_by,
        auth.get_current_identity().to_bytes()
    )

    self.assertEqual(build.bucket_id, 'chromium/try')
    self.assertEqual(build.parameters[model.BUILDER_PARAMETER], 'linux')
    self.assertEqual(build.created_by, auth.get_current_identity())

  def test_add_with_properties(self):
    props = {'foo': 'bar', 'qux': 1}
    prop_struct = bbutil.dict_to_struct(props)
    build = self.add(dict(properties=prop_struct))
    self.assertEqual(build.proto.input.properties, prop_struct)
    self.assertEqual(test_util.msg_to_dict(build.input_properties), props)
    self.assertEqual(build.parameters.get(model.PROPERTIES_PARAMETER), props)

  def test_add_with_notify(self):
    build = self.add(
        dict(
            notify=dict(
                pubsub_topic='projects/p/topics/t',
                user_data='hello',
            )
        ),
    )
    self.assertEqual(build.pubsub_callback.topic, 'projects/p/topics/t')
    self.assertEqual(build.pubsub_callback.user_data, 'hello')

  def test_add_with_gitiles_commit(self):
    gitiles_commit = common_pb2.GitilesCommit(
        host='chromium.googlesource.com',
        project='infra/luci/luci-go',
        ref='refs/heads/master',
        id='b7a757f457487cd5cfe2dae83f65c5bc10e288b7',
        position=1,
    )

    build = self.add(dict(gitiles_commit=gitiles_commit))
    self.assertEqual(build.proto.input.gitiles_commit, gitiles_commit)
    self.assertEqual(build.input_gitiles_commit, gitiles_commit)

  def test_add_with_priority(self):
    build = self.add(dict(priority=42))
    self.assertEqual(build.proto.infra.swarming.priority, 42)

  def test_add_update_builders(self):
    recently = self.now - datetime.timedelta(minutes=1)
    while_ago = self.now - datetime.timedelta(minutes=61)
    ndb.put_multi([
        model.Builder(id='chromium:try:linux', last_scheduled=recently),
        model.Builder(id='chromium:try:mac', last_scheduled=while_ago),
    ])

    creation.add_many_async([
        self.build_request(dict(builder=dict(builder='linux'))),
        self.build_request(dict(builder=dict(builder='mac'))),
        self.build_request(dict(builder=dict(builder='win'))),
    ]).get_result()

    builders = model.Builder.query().fetch()
    self.assertEqual(len(builders), 3)
    self.assertEqual(builders[0].key.id(), 'chromium:try:linux')
    self.assertEqual(builders[0].last_scheduled, recently)
    self.assertEqual(builders[1].key.id(), 'chromium:try:mac')
    self.assertEqual(builders[1].last_scheduled, self.now)
    self.assertEqual(builders[2].key.id(), 'chromium:try:win')
    self.assertEqual(builders[2].last_scheduled, self.now)

  def test_add_with_request_id(self):
    build = self.add(dict(request_id='1'))
    build2 = self.add(dict(request_id='1'))
    self.assertIsNotNone(build.key)
    self.assertEqual(build, build2)

  def test_add_with_leasing(self):
    build = self.add(
        lease_expiration_date=utils.utcnow() + datetime.timedelta(seconds=10),
    )
    self.assertTrue(build.is_leased)
    self.assertGreater(build.lease_expiration_date, utils.utcnow())
    self.assertIsNotNone(build.lease_key)

  def test_add_with_swarming_400(self):
    swarming.create_task_async.return_value = future_exception(
        net.Error('', status_code=400, response='bad request')
    )
    with self.assertRaises(errors.InvalidInputError):
      self.add()

  def test_add_with_build_numbers(self):
    build_numbers = {}

    def create_task_async(build):
      build_numbers[build.parameters['i']] = build.proto.number
      return future(None)

    swarming.create_task_async.side_effect = create_task_async

    (_, ex0), (_, ex1) = creation.add_many_async([
        self.build_request(parameters={'i': 1},),
        self.build_request(parameters={'i': 2},),
    ]).get_result()

    self.assertIsNone(ex0)
    self.assertIsNone(ex1)
    self.assertEqual(build_numbers, {1: 1, 2: 2})

  @mock.patch('sequence.try_return_async', autospec=True)
  def test_add_with_build_numbers_and_return(self, try_return_async):
    try_return_async.return_value = future(None)

    class Error(Exception):
      pass

    swarming.create_task_async.return_value = future_exception(Error())

    with self.assertRaises(Error):
      self.add()

    try_return_async.assert_called_with('chromium/try/linux', 1)

  def test_add_with_swarming_200_and_400(self):

    def create_task_async(b):
      if b.parameters['i'] == 1:
        return future_exception(
            net.Error('', status_code=400, response='bad request')
        )
      b.swarming_hostname = self.chromium_try.swarming.hostname
      b.swarming_task_id = 'deadbeef'
      return future(None)

    swarming.create_task_async.side_effect = create_task_async

    (b0, ex0), (b1, ex1) = creation.add_many_async([
        self.build_request(parameters={'i': 0}),
        self.build_request(parameters={'i': 1})
    ]).get_result()

    self.assertIsNone(ex0)
    self.assertEqual(b0.bucket_id, 'chromium/try')

    self.assertIsNotNone(ex1)
    self.assertIsNone(b1)

  def test_add_with_swarming_403(self):

    swarming.create_task_async.return_value = future_exception(
        net.AuthError('', status_code=403, response='no no')
    )
    with self.assertRaisesRegexp(auth.AuthorizationError, 'no no'):
      self.add()

  def test_builder_tag(self):
    build = self.add(dict(builder=dict(builder='linux')))
    self.assertTrue('builder:linux' in build.tags)

  def test_add_builder_tag_coincide(self):
    build = self.add(
        dict(
            builder=dict(builder='linux'),
            tags=[dict(key='builder', value='linux')],
        )
    )
    self.assertEqual(build.tags, ['builder:linux'])

  def test_buildset_index(self):
    build = self.add(
        dict(
            tags=[
                dict(key='buildset', value='foo'),
                dict(key='buildset', value='bar'),
            ]
        )
    )

    for t in ('buildset:foo', 'buildset:bar'):
      index = search.TagIndex.get_by_id(t)
      self.assertIsNotNone(index)
      self.assertEqual(len(index.entries), 1)
      self.assertEqual(index.entries[0].build_id, build.key.id())
      self.assertEqual(index.entries[0].bucket_id, build.bucket_id)

  def test_buildset_index_with_request_id(self):
    build = self.add(
        dict(
            tags=[dict(key='buildset', value='foo')],
            request_id='0',
        )
    )

    index = search.TagIndex.get_by_id('buildset:foo')
    self.assertIsNotNone(index)
    self.assertEqual(len(index.entries), 1)
    self.assertEqual(index.entries[0].build_id, build.key.id())
    self.assertEqual(index.entries[0].bucket_id, build.bucket_id)

  def test_buildset_index_existing(self):
    search.TagIndex(
        id='buildset:foo',
        entries=[
            search.TagIndexEntry(
                build_id=int(2**63 - 1),
                bucket_id='chromium/try',
            ),
            search.TagIndexEntry(
                build_id=0,
                bucket_id='chromium/try',
            ),
        ]
    ).put()
    build = self.add(dict(tags=[dict(key='buildset', value='foo')]))
    index = search.TagIndex.get_by_id('buildset:foo')
    self.assertIsNotNone(index)
    self.assertEqual(len(index.entries), 3)
    self.assertIn(build.key.id(), [e.build_id for e in index.entries])
    self.assertIn(build.bucket_id, [e.bucket_id for e in index.entries])

  def test_add_many(self):
    results = creation.add_many_async([
        self.build_request(dict(tags=[dict(key='buildset', value='a')])),
        self.build_request(dict(tags=[dict(key='buildset', value='a')])),
    ]).get_result()
    self.assertEqual(len(results), 2)
    self.assertIsNotNone(results[0][0])
    self.assertIsNone(results[0][1])
    self.assertIsNotNone(results[1][0])
    self.assertIsNone(results[1][1])

    self.assertEqual(
        results, sorted(results, key=lambda (b, _): b.key.id(), reverse=True)
    )
    results.reverse()

    index = search.TagIndex.get_by_id('buildset:a')
    self.assertIsNotNone(index)
    self.assertEqual(len(index.entries), 2)
    self.assertEqual(index.entries[0].build_id, results[1][0].key.id())
    self.assertEqual(index.entries[0].bucket_id, results[1][0].bucket_id)
    self.assertEqual(index.entries[1].build_id, results[0][0].key.id())
    self.assertEqual(index.entries[1].bucket_id, results[0][0].bucket_id)

  def test_add_many_with_request_id(self):
    req1 = self.build_request(
        dict(
            tags=[dict(key='buildset', value='a')],
            request_id='0',
        ),
    )
    req2 = self.build_request(dict(tags=[dict(key='buildset', value='a')]))
    creation.add_async(req1).get_result()
    creation.add_many_async([req1, req2]).get_result()

    # Build for req1 must be added only once.
    idx = search.TagIndex.get_by_id('buildset:a')
    self.assertEqual(len(idx.entries), 2)
    self.assertEqual(idx.entries[0].bucket_id, 'chromium/try')

  @mock.patch('search.add_to_tag_index_async', autospec=True)
  def test_add_with_tag_index_contention(self, add_to_tag_index_async):

    def mock_create_task_async(build):
      build.swarming_hostname = 'swarming.example.com'
      build.swarming_task_id = str(build.proto.number)
      return future(None)

    swarming.create_task_async.side_effect = mock_create_task_async
    add_to_tag_index_async.side_effect = Exception('contention')
    swarming.cancel_task_async.side_effect = [
        future(None), future_exception(Exception())
    ]

    with self.assertRaisesRegexp(Exception, 'contention'):
      creation.add_many_async([
          self.build_request(dict(tags=[dict(key='buildset', value='a')])),
          self.build_request(dict(tags=[dict(key='buildset', value='a')])),
      ]).get_result()

    swarming.cancel_task_async.assert_any_call('swarming.example.com', '1')
    swarming.cancel_task_async.assert_any_call('swarming.example.com', '2')
