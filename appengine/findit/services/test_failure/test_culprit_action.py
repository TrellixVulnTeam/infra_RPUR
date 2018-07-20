# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
"""This module is for actions on identified culprits for test failure.

It provides functions to:
  * Determine if Findit should take actions on a culprit.
"""

from datetime import timedelta
import logging

from google.appengine.ext import ndb

from common.waterfall import failure_type
from libs import time_util
from model.wf_suspected_cl import WfSuspectedCL
from services import git
from waterfall import waterfall_config

_DEFAULT_AUTO_CREATE_REVERT_DAILY_THRESHOLD_TEST = 10
_DEFAULT_AUTO_COMMIT_REVERT_DAILY_THRESHOLD_TEST = 4


def _GetDailyNumberOfRevertedCulprits(limit):
  earliest_time = time_util.GetUTCNow() - timedelta(days=1)
  # TODO(chanli): improve the check for a rare case when two pipelines revert
  # at the same time.
  return WfSuspectedCL.query(
      ndb.AND(WfSuspectedCL.failure_type == failure_type.TEST,
              WfSuspectedCL.revert_created_time >= earliest_time)).count(limit)


def CanAutoCreateRevert(culprit, parameters):
  """Checks if Findit can auto create a revert.

  Args:
    culprit (Basestring): Urlsafe key for the suspected cl.
    parameters (CulpritActionParameters): Parameters to run culprit action
      pipelines.

  Findit can auto create a revert if:
    1. Auto create revert for test is turned on;
    2. The number of reverts in past 24 hours is less than the daily limit;
    3. The culprit is also being suspected by the heuristic analysis.
  """
  heuristic_cls = parameters.heuristic_cls
  if culprit not in heuristic_cls:
    return False

  action_settings = waterfall_config.GetActionSettings()
  # Auto revert has been turned off.
  if not bool(action_settings.get('auto_create_revert_test')):
    return False

  auto_create_revert_daily_threshold_test = action_settings.get(
      'auto_create_revert_daily_threshold_test',
      _DEFAULT_AUTO_CREATE_REVERT_DAILY_THRESHOLD_TEST)
  # Auto revert has exceeded daily limit.
  if _GetDailyNumberOfRevertedCulprits(
      auto_create_revert_daily_threshold_test
  ) >= auto_create_revert_daily_threshold_test:
    logging.info('Auto reverts for test culprits on %s has met daily limit.',
                 time_util.FormatDatetime(time_util.GetUTCNow()))
    return False

  return True


def _GetDailyNumberOfCommits(limit):
  earliest_time = time_util.GetUTCNow() - timedelta(days=1)
  # TODO(chanli): improve the check for a rare case when two pipelines commit
  # at the same time.
  return WfSuspectedCL.query(
      ndb.AND(
          WfSuspectedCL.failure_type == failure_type.TEST,
          WfSuspectedCL.revert_committed_time >= earliest_time)).count(limit)


def CanAutoCommitRevertByFindit():
  """Checks if the revert can be auto committed by Findit.

  The revert can be committed if:
    1. Auto revert and Auto commit is turned on;
    2. The number of commits of reverts in past 24 hours is less than the
      daily limit;
    3. Culprit author has not landed another change yet.
  """
  action_settings = waterfall_config.GetActionSettings()
  if (not bool(action_settings.get('auto_commit_revert_test')) or
      not bool(action_settings.get('auto_create_revert_test'))):
    return False

  auto_commit_revert_daily_threshold_test = action_settings.get(
      'auto_commit_revert_daily_threshold_test',
      _DEFAULT_AUTO_COMMIT_REVERT_DAILY_THRESHOLD_TEST)
  if _GetDailyNumberOfCommits(auto_commit_revert_daily_threshold_test
                             ) >= auto_commit_revert_daily_threshold_test:
    logging.info('Auto commits on %s has met daily limit.',
                 time_util.FormatDatetime(time_util.GetUTCNow()))
    return False
  return True
