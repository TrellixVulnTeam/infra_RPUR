# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
"""Utilities to assist reporting analysis results."""
import textwrap

from google.appengine.ext import ndb
from libs import time_util
from model.flake.analysis import master_flake_analysis
from model.flake.flake import Flake
from services import build_url
from services import issue_tracking_service
from services import monitoring
from services import swarming
from services.flake_failure import flake_constants
from waterfall import waterfall_config

_RESULT_WITH_CULPRIT_TEMPLATE = textwrap.dedent("""
Flaky test: {test_name}
Sample failed build due to flakiness: {build_link}
Test output log: {test_output_log_link}
Culprit ({confidence_score} confidence): r{commit_position}
Analysis: {analysis_link}

Please revert the culprit, or disable the test and find the appropriate owner.

If the culprit above is wrong, please file a bug using this link:
{wrong_result_link}

Automatically posted by the findit-for-me app (https://goo.gl/Ot9f7N).""")

_WRONG_RESULT_LINK_TEMPLATE = (
    'https://bugs.chromium.org/p/chromium/issues/entry?'
    'status=Unconfirmed&'
    'labels=Pri-1,Test-Findit-Wrong&'
    'components=Tools%3ETest%3EFindit%3EFlakiness&'
    'summary=%5BFindit%5D%20Flake%20Analyzer%20-%20Wrong%20result%20for%20{}&'
    'comment=Link%20to%20Analysis%3A%20{}')

_UNKNOWN_CULPRIT_TEMPLATE = textwrap.dedent("""
Flaky test: {test_name}
Sample failed build due to flakiness: {build_link}
Test output log: {test_output_log_link}
Analysis: {analysis_link}

This flake is either longstanding, has low flakiness, or is not reproducible.

Automatically posted by the findit-for-me app (https://goo.gl/Ot9f7N).""")

_FINDIT_ANALYZED_LABEL_TEXT = 'Test-Findit-Analyzed'


def _GenerateAnalysisLink(analysis):
  """Returns a link to Findit's result page of a MasterFlakeAnalysis."""
  return 'https://findit-for-me.appspot.com/waterfall/flake?key={}'.format(
      analysis.key.urlsafe())


def _GenerateWrongResultLink(analysis):
  """Returns the test with a link to file a bug agasinst a wrong result."""
  return _WRONG_RESULT_LINK_TEMPLATE.format(analysis.test_name,
                                            _GenerateAnalysisLink(analysis))


def _GenerateTestOutputLogLink(analysis):
  """Generates a link to the swarming task to be surfaced to the bug.

  Args:
    analysis (MasterFlakeAnalysis): The analysis whose data points and swarming
        tasks will be queried for surfacing to the bug.

  Returns:
    url (str): The url to the swarming task.
  """
  task_id = analysis.GetRepresentativeSwarmingTaskId()
  assert task_id, 'Representative task id unexpectedly not found!'

  return swarming.GetSwarmingTaskUrl(task_id)


def _GenerateMessageText(analysis):
  """Generates the text to create or update a bug with depending on results.

  Args:
    analysis (MasterFlakeAnalysis): The completed analysis with results to
      determine what to update the bug with.

  Returns:
    (str): The text to upodate the bug with.
  """
  analysis_link = _GenerateAnalysisLink(analysis)

  build_link = build_url.CreateBuildUrl(analysis.original_master_name,
                                        analysis.original_builder_name,
                                        analysis.original_build_number)
  test_output_log_link = _GenerateTestOutputLogLink(analysis)

  if analysis.culprit_urlsafe_key:
    culprit = ndb.Key(urlsafe=analysis.culprit_urlsafe_key).get()
    assert culprit, 'Culprit is unexpectedly missing.'

    culprit_confidence = "{0:0.1f}%".format(
        analysis.confidence_in_culprit * 100)
    wrong_result_link = _GenerateWrongResultLink(analysis)

    return _RESULT_WITH_CULPRIT_TEMPLATE.format(
        test_name=analysis.test_name,
        build_link=build_link,
        test_output_log_link=test_output_log_link,
        confidence_score=culprit_confidence,
        commit_position=culprit.commit_position,
        analysis_link=analysis_link,
        wrong_result_link=wrong_result_link)

  return _UNKNOWN_CULPRIT_TEMPLATE.format(
      test_name=analysis.original_test_name,
      build_link=build_link,
      test_output_log_link=test_output_log_link,
      analysis_link=analysis_link)


class FlakeAnalysisIssueGenerator(
    issue_tracking_service.FlakyTestIssueGenerator):
  """Encapsulates the details of issues filed by Flake Analyzer."""

  def __init__(self, analysis):
    super(FlakeAnalysisIssueGenerator, self).__init__()
    self._analysis = analysis

  def GetStepName(self):
    return Flake.NormalizeStepName(
        step_name=self._analysis.step_name,
        master_name=self._analysis.master_name,
        builder_name=self._analysis.builder_name,
        build_number=self._analysis.build_number)

  def GetTestName(self):
    return Flake.NormalizeTestName(self._analysis.test_name,
                                   self._analysis.step_name)

  def GetTestLabelName(self):
    # Issues are filed with the test label name.
    return Flake.GetTestLabelName(self._analysis.test_name,
                                  self._analysis.step_name)

  def GetMonorailProject(self):
    # Currently, flake analysis only works on Chromium project.
    return 'chromium'

  def GetDescription(self):
    return _GenerateMessageText(self._analysis)

  def GetComment(self):
    return _GenerateMessageText(self._analysis)

  def ShouldRestoreChromiumSheriffLabel(self):
    # Analysis results are not always immediately actionable, so don't restore
    # Sheriff label to avoid being too noisy.
    return False

  def GetLabels(self):
    priority = self.GetPriority()
    flaky_test_labels = self._GetCommonFlakyTestLabel()
    flaky_test_labels.append(priority)
    flaky_test_labels.append(_FINDIT_ANALYZED_LABEL_TEXT)
    return flaky_test_labels

  def OnIssueCreated(self):
    monitoring.OnIssueChange('created', 'flake')

  def OnIssueUpdated(self):
    monitoring.OnIssueChange('update', 'flake')


def ShouldFileBugForAnalysis(analysis):
  """Returns true if a bug should be filed for this analysis.

  The requirements for a bug to be filed.
    - The bug creation feature if enabled.
    - The pipeline hasn't been attempted before (see above).
    - The analysis has sufficient confidence (1.0).
    - The analysis doesn't already have a bug associated with it.
    - A duplicate bug hasn't been filed by Findit or CTF.
    - A duplicate bug hasn't been filed by a human.
  """
  if not UnderDailyLimit():
    analysis.LogInfo('Reached bug filing limit for the day.')
    return False

  if HasPreviousAttempt(analysis):
    analysis.LogWarning(
        'There has already been an attempt at filing a bug, aborting.')
    return False

  if not HasSufficientConfidenceInCulprit(analysis,
                                          GetMinimumConfidenceToFileBugs()):
    analysis.LogInfo('''Analysis has confidence {:.2%}
        which isn\'t high enough to file a bug.'''.format(
        analysis.confidence_in_culprit))
    return False

  # Check if there's already a bug attached to this issue.
  if analysis.bug_id and issue_tracking_service.OpenBugAlreadyExistsForId(
      analysis.bug_id):
    analysis.LogInfo('Bug with id {} already exists.'.format(analysis.bug_id))
    return False

  if issue_tracking_service.OpenIssueAlreadyExistsForFlakyTest(
      analysis.test_name):
    analysis.LogInfo('Bug already exists for flaky test: {}'.format(
        analysis.test_name))
    return False

  return True


def ShouldUpdateBugForAnalysis(analysis):
  assert not analysis.error

  if not analysis.bug_id:
    analysis.LogInfo('bug=%s' % analysis.bug_id)
    if analysis.culprit_urlsafe_key:
      monitoring.OnFlakeCulprit('culprit-identified', 'none',
                                'no-bug-to-update')
    else:
      monitoring.OnFlakeCulprit('culprit-not-identified', 'none',
                                'no-bug-to-update')
    return False

  if len(analysis.data_points) < 2:
    analysis.LogInfo('%d data points' % len(analysis.data_points))
    monitoring.OnFlakeCulprit('culprit-identified', 'none',
                              'insufficient-datapoints')
    return False

  if (analysis.culprit_urlsafe_key and not HasSufficientConfidenceInCulprit(
      analysis, GetMinimumConfidenceToUpdateBugs())):
    # There is a culprit, but insufficient confidence.
    monitoring.OnFlakeCulprit('culprit-identified', 'none',
                              'insufficient-confidence')
    return False

  # TODO(crbug.com/847960): Do not update bugs if Findit already logged one as
  # the information would be redundant.
  return True


def UnderDailyLimit():
  action_settings = waterfall_config.GetActionSettings()
  daily_bug_limit = action_settings.get(
      'max_flake_bug_updates_per_day',
      flake_constants.DEFAULT_MAX_BUG_UPDATES_PER_DAY)
  query = master_flake_analysis.MasterFlakeAnalysis.query(
      master_flake_analysis.MasterFlakeAnalysis.request_time >= time_util
      .GetMostRecentUTCMidnight())
  bugs_filed_today = 0

  more = True
  cursor = None
  while more:
    results, cursor, more = query.fetch_page(100, start_cursor=cursor)
    for result in results:
      if result.has_attempted_filing and result.bug_id:
        bugs_filed_today += 1

  return bugs_filed_today < daily_bug_limit


def HasSufficientConfidenceInCulprit(analysis, required_confidence):
  """Returns true is there's high enough confidence in the culprit."""
  if not analysis.confidence_in_culprit:
    return False

  return (analysis.confidence_in_culprit + flake_constants.EPSILON >=
          required_confidence)


def GetMinimumConfidenceToUpdateBugs():
  return waterfall_config.GetCheckFlakeSettings().get(
      'minimum_confidence_to_update_cr',
      flake_constants.DEFAULT_MINIMUM_CONFIDENCE_SCORE_TO_UPDATE_CR)


def GetMinimumConfidenceToFileBugs():
  return waterfall_config.GetCheckFlakeSettings().get(
      'minimum_confidence_to_create_bug',
      flake_constants.DEFAULT_MINIMUM_CONFIDENCE_TO_CREATE_BUG)


def HasPreviousAttempt(analysis):
  """Returns True if an analysis has already attempted to file a bug."""
  return analysis.has_attempted_filing
