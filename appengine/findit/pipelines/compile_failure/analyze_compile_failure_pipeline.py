# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import logging

from common import constants
from common import monitoring
from dto.start_waterfall_try_job_inputs import StartCompileTryJobInput
from gae_libs import appengine_util
from gae_libs import pipelines
from gae_libs.pipeline_wrapper import BasePipeline
from libs import analysis_status
from libs import time_util
from model.wf_analysis import WfAnalysis
from pipelines import report_event_pipeline
from pipelines.compile_failure.heuristic_analysis_for_compile_pipeline import (
    HeuristicAnalysisForCompilePipeline)
from pipelines.compile_failure.start_compile_try_job_pipeline import (
    StartCompileTryJobPipeline)
from services.parameters import BuildKey
from services.parameters import CompileFailureInfo
from services.parameters import CompileHeuristicAnalysisOutput
from services.parameters import CompileHeuristicAnalysisParameters


class AnalyzeCompileFailurePipeline(BasePipeline):

  def __init__(self, master_name, builder_name, build_number,
               current_failure_info, build_completed, force):
    super(AnalyzeCompileFailurePipeline,
          self).__init__(master_name, builder_name, build_number,
                         current_failure_info, build_completed, force)
    self.master_name = master_name
    self.builder_name = builder_name
    self.build_number = build_number
    self.current_failure_info = current_failure_info
    self.build_completed = build_completed
    self.force = force

  def _HandleUnexpectedAborting(self, was_aborted):
    """Handles unexpected aborting gracefully.

    Marks the WfAnalysis status as error, indicating that it was aborted.
    If one of heuristic pipelines caused the abort, continue try job analysis
    by starting a new pipeline.

    Args:
      was_aborted (bool): True if the pipeline was aborted, otherwise False.
    """
    if not was_aborted:
      return

    analysis = WfAnalysis.Get(self.master_name, self.builder_name,
                              self.build_number)
    # Heuristic analysis could have already completed, while triggering the
    # try job kept failing and lead to the abort.
    run_try_job = False
    if not analysis.completed:
      # Heuristic analysis is aborted.
      analysis.status = analysis_status.ERROR
      analysis.result_status = None

      if analysis.failure_info:
        # We need failure_info to run try jobs,
        # while signals is optional for compile try jobs.
        run_try_job = True
    analysis.aborted = True
    analysis.put()

    monitoring.aborted_pipelines.increment({'type': 'compile'})

    if not run_try_job:
      return

    self._ContinueTryJobPipeline(analysis.failure_info, analysis.signals)

  def finalized(self):
    self._HandleUnexpectedAborting(self.was_aborted)
    monitoring.completed_pipelines.increment({'type': 'compile'})

  def _ContinueTryJobPipeline(self, failure_info, signals):
    heuristic_result = {
        'failure_info': failure_info,
        'signals': signals,
        'heuristic_result': None
    }
    start_compile_try_job_input = StartCompileTryJobInput(
        build_key=BuildKey(
            master_name=self.master_name,
            builder_name=self.builder_name,
            build_number=self.build_number),
        heuristic_result=CompileHeuristicAnalysisOutput.FromSerializable(
            heuristic_result),
        build_completed=self.build_completed,
        force=self.force)
    try_job_pipeline = StartCompileTryJobPipeline(start_compile_try_job_input)
    try_job_pipeline.target = appengine_util.GetTargetNameForModule(
        constants.WATERFALL_BACKEND)
    try_job_pipeline.start(queue_name=constants.WATERFALL_ANALYSIS_QUEUE)
    logging.info(
        'A try job pipeline for build %s, %s, %s starts after heuristic '
        'analysis was aborted. Check pipeline at: %s.', self.master_name,
        self.builder_name, self.build_number, self.pipeline_status_path())

  def _ResetAnalysis(self, master_name, builder_name, build_number):
    analysis = WfAnalysis.Get(master_name, builder_name, build_number)
    analysis.pipeline_status_path = self.pipeline_status_path()
    analysis.status = analysis_status.RUNNING
    analysis.result_status = None
    analysis.start_time = time_util.GetUTCNow()
    analysis.version = appengine_util.GetCurrentVersion()
    analysis.end_time = None
    analysis.put()

  # Arguments number differs from overridden method - pylint: disable=W0221
  def run(self, master_name, builder_name, build_number, current_failure_info,
          build_completed, force):
    self._ResetAnalysis(master_name, builder_name, build_number)

    # The yield statements below return PipelineFutures, which allow subsequent
    # pipelines to refer to previous output values.
    # https://github.com/GoogleCloudPlatform/appengine-pipelines/wiki/Python

    # Heuristic Approach.
    heuristic_params = CompileHeuristicAnalysisParameters(
        failure_info=CompileFailureInfo.FromSerializable(current_failure_info),
        build_completed=build_completed)
    heuristic_result = yield HeuristicAnalysisForCompilePipeline(
        heuristic_params)

    # Try job approach.
    # Checks if first time failures happen and starts a try job if yes.
    with pipelines.pipeline.InOrder():
      start_compile_try_job_input = pipelines.CreateInputObjectInstance(
          StartCompileTryJobInput,
          build_key=BuildKey(
              master_name=master_name,
              builder_name=builder_name,
              build_number=build_number),
          heuristic_result=heuristic_result,
          build_completed=build_completed,
          force=force)
      yield StartCompileTryJobPipeline(start_compile_try_job_input)
      # Report event to BQ.
      report_event_input = pipelines.CreateInputObjectInstance(
          report_event_pipeline.ReportEventInput,
          analysis_urlsafe_key=WfAnalysis.Get(master_name, builder_name,
                                              build_number).key.urlsafe())
      if not force:
        yield report_event_pipeline.ReportAnalysisEventPipeline(
            report_event_input)
