# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import timedelta
import logging
import random

from common import appengine_util
from common import constants
from common.pipeline_wrapper import BasePipeline
from common.pipeline_wrapper import pipeline
from libs import time_util
from model import analysis_status
from model import result_status
from model.flake.flake_swarming_task import FlakeSwarmingTask
from model.flake.master_flake_analysis import MasterFlakeAnalysis
from waterfall import waterfall_config
from waterfall.flake import confidence
from waterfall.flake import lookback_algorithm
from waterfall.flake.lookback_algorithm import NormalizedDataPoint
from waterfall.flake.recursive_flake_try_job_pipeline import CreateCulprit
from waterfall.flake.recursive_flake_try_job_pipeline import (
    RecursiveFlakeTryJobPipeline)
from waterfall.flake.recursive_flake_try_job_pipeline import (
    UpdateAnalysisTryJobStatusUponCompletion)
from waterfall.flake.update_flake_bug_pipeline import UpdateFlakeBugPipeline
from waterfall.process_flake_swarming_task_result_pipeline import (
    ProcessFlakeSwarmingTaskResultPipeline)
from waterfall.trigger_flake_swarming_task_pipeline import (
    TriggerFlakeSwarmingTaskPipeline)


_DEFAULT_MINIMUM_CONFIDENCE_SCORE = 0.6
_DEFAULT_MAX_BUILD_NUMBERS = 500


def _UpdateAnalysisStatusUponCompletion(
    analysis, suspected_build, status, error, build_confidence_score=None,
    try_job_status=analysis_status.SKIPPED):
  if suspected_build is None:
    analysis.end_time = time_util.GetUTCNow()
    analysis.result_status = result_status.NOT_FOUND_UNTRIAGED
  else:
    analysis.suspected_flake_build_number = suspected_build

  analysis.error = error
  analysis.status = status
  analysis.try_job_status = try_job_status
  analysis.confidence_in_suspected_build = build_confidence_score
  analysis.put()


def _GetETAToStartAnalysis(manually_triggered):
  """Returns an ETA as of a UTC datetime.datetime to start the analysis.

  If not urgent, Swarming tasks should be run off PST peak hours from 11am to
  6pm on workdays.

  Args:
    manually_triggered (bool): True if the analysis is from manual request, like
        by a Chromium sheriff.

  Returns:
    The ETA as of a UTC datetime.datetime to start the analysis.
  """
  if manually_triggered:
    # If the analysis is manually triggered, run it right away.
    return time_util.GetUTCNow()

  now_at_pst = time_util.GetDatetimeInTimezone(
      'US/Pacific', time_util.GetUTCNowWithTimezone())
  if now_at_pst.weekday() >= 5:  # PST Saturday or Sunday.
    return time_util.GetUTCNow()

  if now_at_pst.hour < 11 or now_at_pst.hour >= 18:  # Before 11am or after 6pm.
    return time_util.GetUTCNow()

  # Set ETA time to 6pm, and also with a random latency within 30 minutes to
  # avoid sudden burst traffic to Swarming.
  diff = timedelta(hours=18 - now_at_pst.hour,
                   minutes=-now_at_pst.minute,
                   seconds=-now_at_pst.second + random.randint(0, 30 * 60),
                   microseconds=-now_at_pst.microsecond)
  eta = now_at_pst + diff

  # Convert back to UTC.
  return time_util.GetDatetimeInTimezone('UTC', eta)


def _IsSwarmingTaskSufficientForCacheHit(
    flake_swarming_task, number_of_iterations):
  """Determines whether or not a swarming task is sufficient for a cache hit.

  Args:
    flake_swarming_task (FlakeSwarmingTask): The task to be examined.
    number_of_iterations (int): The minimum number of iterations
      flake_swarming_task needs to have run in order to count as a cache hit.

  Returns:
    A bool whether or not flake_swarming_task is sufficient to be a cache hit.
  """
  # Swarming task must exist.
  if not flake_swarming_task:
    return False

  # Cached swarming task's numbers must be thorough enough.
  if flake_swarming_task.tries < number_of_iterations:
    return False

  # Cached swarming task must either be scheduled, in progress, or completed.
  return flake_swarming_task.status in [analysis_status.PENDING,
                                        analysis_status.RUNNING,
                                        analysis_status.COMPLETED]


def _GetBestBuildNumberToRun(
    master_name, builder_name, preferred_run_build_number, step_name, test_name,
    step_size, number_of_iterations):
  """Finds the optimal nearby swarming task build number to use for a cache hit.

  Builds are searched back looking for something either already completed or in
  progress. Completed builds are returned immediately, whereas for those in
  progress the closer the build number is to the original, the higher priority
  it is given.

  Args:
    master_name (str): The name of the master for this flake analysis.
    builder_name (str): The name of the builder for this flake analysis.
    preferred_run_build_number (int): The originally-requested build number to
      run the swarming task on.
    step_name (str): The name of the step to run swarming on.
    test_name (str): The name of the test to run swarming on.
    step_size (int): The distance of the last preferred build number that was
      called on this analysis. Used for determining the lookback threshold.
    number_of_iterations (int): The number of iterations being requested for
      the swarming task that is to be performed. Used to determine a sufficient
      cache hit.

  Returns:
    build_number (int): The best build number to analyze for this iteration of
      the flake analysis.
  """
  # Looks forward or backward up to half of step_size.
  possibly_cached_build_numbers = _GetListOfNearbyBuildNumbers(
      preferred_run_build_number, step_size / 2)
  candidate_build_number = None
  candidate_flake_swarming_task_status = None

  for build_number in possibly_cached_build_numbers:
    cached_flake_swarming_task = FlakeSwarmingTask.Get(
        master_name, builder_name, build_number, step_name, test_name)
    sufficient = _IsSwarmingTaskSufficientForCacheHit(
        cached_flake_swarming_task, number_of_iterations)

    if sufficient:
      if cached_flake_swarming_task.status == analysis_status.COMPLETED:
        # Found a nearby swarming task that's already done.
        return build_number

      # Keep searching, but keeping this candidate in mind. Pending tasks are
      # considered, but running tasks are given higher priority.
      # TODO(lijeffrey): A further optimization can be to pick the swarming
      # task with the earliest ETA.
      if (candidate_build_number is None or
          (candidate_flake_swarming_task_status == analysis_status.PENDING and
           cached_flake_swarming_task.status == analysis_status.RUNNING)):
        # Either no previous candidate or a better candidate was found.
        candidate_build_number = build_number
        candidate_flake_swarming_task_status = cached_flake_swarming_task.status

  # No cached build nearby deemed adequate could be found.
  return candidate_build_number or preferred_run_build_number


def _GetListOfNearbyBuildNumbers(preferred_run_build_number, maximum_threshold):
  """Gets a list of numbers within range near preferred_run_build_number.

  Args:
    preferred_run_build_number (int): Assumed to be a positive number.
    maximum_threshold (int): A non-negative number for how far in either
    direction to look.

  Returns:
    A list of nearby numbers within maximum_threshold before and after
    preferred_run_build_number, ordered by closest to farthest. For example, if
    preferred_run_build_number is 1000 and maximum_threshold is 2, return
    [1000, 999, 1001, 998, 1002].
  """
  if maximum_threshold >= preferred_run_build_number:
    # Build numbers are always assumed to start from 1, so don't include
    # anything before that.
    return range(1, preferred_run_build_number + maximum_threshold + 1)

  nearby_build_numbers = [preferred_run_build_number]

  for i in range(1, maximum_threshold + 1):
    nearby_build_numbers.append(preferred_run_build_number - i)
    nearby_build_numbers.append(preferred_run_build_number + i)

  return nearby_build_numbers


class RecursiveFlakePipeline(BasePipeline):

  def __init__(self, *args, **kwargs):
    super(RecursiveFlakePipeline, self).__init__(*args, **kwargs)
    self.manually_triggered = kwargs.get('manually_triggered', False)

  def StartOffPSTPeakHours(self, *args, **kwargs):
    """Starts the pipeline off PST peak hours if not triggered manually."""
    kwargs['eta'] = _GetETAToStartAnalysis(self.manually_triggered)
    self.start(*args, **kwargs)

  # Arguments number differs from overridden method - pylint: disable=W0221
  def run(self, master_name, builder_name, preferred_run_build_number,
          step_name, test_name, version_number, triggering_build_number,
          manually_triggered=False, use_nearby_neighbor=False, step_size=0):
    """Pipeline to determine the regression range of a flaky test.

    Args:
      master_name (str): The master name.
      builder_name (str): The builder name.
      preferred_run_build_number (int): The build number the check flake
        algorithm should perform a swarming rerun on, but may be overridden to
        use the results of a nearby neighbor if use_nearby_neighbor is True.
      step_name (str): The step name.
      test_name (str): The test name.
      version_number (int): The version to save analysis results and data to.
      triggering_build_number (int): The build number that triggered this
        analysis.
      manually_triggered (bool): True if the analysis is from manual request,
        like by a Chromium sheriff.
      use_nearby_neighbor (bool): Whether the optimization for using the
        swarming results of a nearby build number, if available, should be used
        in place of triggering a new swarming task on
        preferred_run_build_number.
      step_size (int): The difference in build numbers since the last call to
        RecursiveFlakePipeline to determine the bounds for how far a nearby
        build's swarming task results should be used. Only relevant if
        use_nearby_neighbor is True.
    Returns:
      A dict of lists for reliable/flaky tests.
    """
    flake_analysis = MasterFlakeAnalysis.GetVersion(
        master_name, builder_name, triggering_build_number, step_name,
        test_name, version=version_number)
    logging.info(
        'Running RecursiveFlakePipeline on MasterFlakeAnalysis %s/%s/%s/%s/%s',
        master_name, builder_name, triggering_build_number, step_name,
        test_name)
    logging.info(
        'MasterFlakeAnalysis %s version %s', flake_analysis, version_number)

    if flake_analysis.status != analysis_status.RUNNING:  # pragma: no branch
      flake_analysis.status = analysis_status.RUNNING
      flake_analysis.start_time = time_util.GetUTCNow()
      flake_analysis.put()

    # TODO(lijeffrey): Allow custom parameters supplied by user.
    iterations = waterfall_config.GetCheckFlakeSettings().get(
        'swarming_rerun', {}).get('iterations_to_rerun', 100)
    actual_run_build_number = _GetBestBuildNumberToRun(
        master_name, builder_name, preferred_run_build_number, step_name,
        test_name, step_size, iterations) if use_nearby_neighbor else (
            preferred_run_build_number)

    # Call trigger pipeline (flake style).
    task_id = yield TriggerFlakeSwarmingTaskPipeline(
        master_name, builder_name, actual_run_build_number, step_name,
        [test_name])

    with pipeline.InOrder():
      yield ProcessFlakeSwarmingTaskResultPipeline(
          master_name, builder_name, actual_run_build_number, step_name,
          task_id, triggering_build_number, test_name, version_number)
      yield NextBuildNumberPipeline(
          master_name, builder_name, triggering_build_number,
          actual_run_build_number, step_name, test_name, version_number,
          use_nearby_neighbor=use_nearby_neighbor,
          manually_triggered=manually_triggered)


def _NormalizeDataPoints(data_points):
  normalized_data_points = [
      (lambda data_point: NormalizedDataPoint(
          data_point.build_number, data_point.pass_rate))(
              d) for d in data_points]
  return sorted(
      normalized_data_points, key=lambda k: k.run_point_number, reverse=True)


class NextBuildNumberPipeline(BasePipeline):

  # Arguments number differs from overridden method - pylint: disable=W0221
  # Unused argument - pylint: disable=W0613
  def run(
      self, master_name, builder_name, triggering_build_number,
      current_build_number, step_name, test_name, version_number,
      use_nearby_neighbor=False, manually_triggered=False):
    # Get MasterFlakeAnalysis success list corresponding to parameters.
    analysis = MasterFlakeAnalysis.GetVersion(
        master_name, builder_name, triggering_build_number, step_name,
        test_name, version=version_number)

    flake_swarming_task = FlakeSwarmingTask.Get(
        master_name, builder_name, current_build_number, step_name, test_name)

    # Don't call another pipeline if we fail.
    if flake_swarming_task.status == analysis_status.ERROR:
      # Report the last flake swarming task's error that it encountered.
      # TODO(lijeffrey): Another neighboring swarming task may be needed in this
      # one's place instead of failing altogether.
      error = flake_swarming_task.error or {
          'error': 'Swarming task failed',
          'message': 'The last swarming task did not complete as expected'
      }

      _UpdateAnalysisStatusUponCompletion(
          analysis, None, analysis_status.ERROR, error)
      logging.error('Error in Swarming task')
      yield UpdateFlakeBugPipeline(analysis.key.urlsafe())
      return

    flake_settings = waterfall_config.GetCheckFlakeSettings()
    algorithm_settings = flake_settings.get('swarming_rerun', {})

    data_points = _NormalizeDataPoints(analysis.data_points)
    # Figure out what build_number to trigger a swarming rerun on next, if any.
    (next_build_number,
     suspected_build) = lookback_algorithm.GetNextRunPointNumber(
         data_points, algorithm_settings)

    max_build_numbers_to_look_back = algorithm_settings.get(
        'max_build_numbers_to_look_back', _DEFAULT_MAX_BUILD_NUMBERS)
    last_build_number = max(
        0, triggering_build_number - max_build_numbers_to_look_back)

    if (next_build_number < last_build_number or
        next_build_number >= triggering_build_number):  # Finished.
      build_confidence_score = None
      if suspected_build is not None:
        # Use steppiness as the confidence score.
        build_confidence_score = confidence.SteppinessForBuild(
            analysis.data_points, suspected_build)

      # Update suspected build and the confidence score.
      _UpdateAnalysisStatusUponCompletion(
          analysis, suspected_build, analysis_status.COMPLETED,
          None, build_confidence_score=build_confidence_score)

      minimum_confidence_score_to_run_tryjobs = flake_settings.get(
          'minimum_confidence_score_to_run_tryjobs',
          _DEFAULT_MINIMUM_CONFIDENCE_SCORE)

      if build_confidence_score is None:
        logging.info(('Skipping try jobs due to no suspected flake build being '
                      'identified'))
      elif build_confidence_score < minimum_confidence_score_to_run_tryjobs:
        # If confidence is too low, bail out on try jobs. Based on analysis of
        # historical data, 60% confidence could filter out almost all false
        # positives.
        analysis.result_status = result_status.FOUND_UNTRIAGED
        analysis.put()
      else:
        # Hook up with try-jobs.
        suspected_build_point = analysis.GetDataPointOfSuspectedBuild()
        assert suspected_build_point

        if suspected_build_point.blame_list:
          if len(suspected_build_point.blame_list) > 1:
            logging.info('Running try-jobs against commits in suspected build')
            start_commit_position = suspected_build_point.commit_position - 1
            start_revision = suspected_build_point.GetRevisionAtCommitPosition(
                start_commit_position)
            yield RecursiveFlakeTryJobPipeline(
                analysis.key.urlsafe(), start_commit_position, start_revision)
            return  # No update to bug yet.
          else:
            logging.info('Single commit in the blame list of suspected build')
            culprit_confidence_score = confidence.SteppinessForCommitPosition(
                analysis.data_points, suspected_build_point.commit_position)
            culprit = CreateCulprit(
                suspected_build_point.git_hash,
                suspected_build_point.commit_position,
                culprit_confidence_score)
            UpdateAnalysisTryJobStatusUponCompletion(
                analysis, culprit, analysis_status.COMPLETED, None)
        else:
          logging.error('Cannot run flake try jobs against empty blame list')
          error = {
              'error': 'Could not start try jobs',
              'message': 'Empty blame list'
          }
          UpdateAnalysisTryJobStatusUponCompletion(
              analysis, None, analysis_status.ERROR, error)

      yield UpdateFlakeBugPipeline(analysis.key.urlsafe())
      return

    pipeline_job = RecursiveFlakePipeline(
        master_name, builder_name, next_build_number, step_name, test_name,
        version_number, triggering_build_number,
        manually_triggered=manually_triggered,
        use_nearby_neighbor=use_nearby_neighbor,
        step_size=(current_build_number - next_build_number))
    # Disable attribute 'target' defined outside __init__ pylint warning,
    # because pipeline generates its own __init__ based on run function.
    pipeline_job.target = (  # pylint: disable=W0201
        appengine_util.GetTargetNameForModule(constants.WATERFALL_BACKEND))
    pipeline_job.StartOffPSTPeakHours(
        queue_name=self.queue_name or constants.DEFAULT_QUEUE)
