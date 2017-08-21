# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from collections import OrderedDict
import mock

from analysis.type_enums import CrashClient
from common.appengine_testcase import AppengineTestCase
from common.model.clusterfuzz_analysis import ClusterfuzzAnalysis
from frontend.handlers import clusterfuzz_dashboard
from frontend.handlers.clusterfuzz_dashboard import GetCommitsInRegressionRange


class ClusterfuzzDashBoardTest(AppengineTestCase):

  def setUp(self):
    super(ClusterfuzzDashBoardTest, self).setUp()
    self.dashboard = clusterfuzz_dashboard.ClusterfuzzDashBoard()

  def testGetCommitsInRegressionRange(self):

    regression_range = {'repo_url': 'http://repo',
                        'repo_path': 'src',
                        'old_revision': 'rev2',
                        'new_revision': 'rev7'}

    class MockRepository(object):

      def GetCommitsBetweenRevisions(self, old_revision, new_revision):
        assert old_revision == regression_range['old_revision']
        assert new_revision == regression_range['new_revision']
        return ['rev1', 'rev2', 'rev3']

    def MockGetRepository(repo_url):
      assert repo_url == regression_range['repo_url']
      return MockRepository()

    self.assertEqual(GetCommitsInRegressionRange(regression_range,
                                                 MockGetRepository), 3)

  def testCrashAnalysisCls(self):
    self.assertEqual(self.dashboard.crash_analysis_cls, ClusterfuzzAnalysis)

  def testClient(self):
    self.assertEqual(self.dashboard.client, CrashClient.CLUSTERFUZZ)

  def testTemplate(self):
    self.assertEqual(self.dashboard.template, 'clusterfuzz_dashboard.html')

  def testPropertyToValueConverter(self):
    self.assertListEqual(
        list(self.dashboard.property_to_value_converter),
        ['found_suspects',
        'has_regression_range',
        'suspected_cls_triage_status',
        'regression_range_triage_status',
        'testcase_id'])
    self.assertTrue(self.dashboard.property_to_value_converter[
        'found_suspects']('yes'))
    self.assertTrue(self.dashboard.property_to_value_converter[
        'has_regression_range']('yes'))
    self.assertEqual(self.dashboard.property_to_value_converter[
        'suspected_cls_triage_status']('1'), 1)
    self.assertEqual(self.dashboard.property_to_value_converter[
        'regression_range_triage_status']('2'), 2)
    self.assertEqual(self.dashboard.property_to_value_converter[
        'testcase_id']('123'), '123')

  def testCrashDataToDisplayWhenThereIsNoCrashToDisplay(self):
    self.assertEqual(self.dashboard.CrashDataToDisplay([]), [])

  @mock.patch('frontend.handlers.clusterfuzz_dashboard.'
              'GetCommitsInRegressionRange')
  def testCrashDataToDisplay(self, get_commits):
    get_commits.return_value = 3

    analysis = ClusterfuzzAnalysis()
    analysis.signature = 'sig'
    analysis.testcase_id = '123'
    analysis.crashed_version = '134abs'
    analysis.job_type = 'asan_job'
    analysis.crash_type = 'check'
    analysis.platform = 'win'
    analysis.sanitizer = 'asan'
    analysis.regression_range = {'repo_path': 'src', 'repo_url': 'https://repo',
                                 'old_revision': 'rev2', 'new_revision': 'rev8'}
    analysis.commits = 3
    analysis.error_name = 'Failed to parse stacktrace'
    analysis.result = {
        'suspected_cls': [{'author': 'someone'}],
        'suspected_project': 'chromium',
        'suspected_components': ['Blink'],
    }
    analysis.put()

    expected_display_data = [{
        'signature': 'sig',
        'testcase_id': '123',
        'version': '134abs',
        'job_type': 'asan_job',
        'crash_type': 'check',
        'platform': 'win',
        'sanitizer': 'asan',
        'regression_range': ['rev2', 'rev8'],
        'commits': 3,
        'error_name': 'Failed to parse stacktrace',
        'suspected_cls': [{'author': 'someone'}],
        'suspected_project': 'chromium',
        'suspected_components': ['Blink'],
        'key': analysis.key.urlsafe(),
    }]

    self.assertListEqual(self.dashboard.CrashDataToDisplay([analysis]),
                         expected_display_data)
