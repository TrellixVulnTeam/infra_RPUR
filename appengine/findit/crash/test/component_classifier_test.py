# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from crash.stacktrace import StackFrame
from crash.stacktrace import CallStack
from crash.component_classifier import Component
from crash.component_classifier import ComponentClassifier
from crash.suspect import Suspect
from crash.test.predator_testcase import PredatorTestCase
from gae_libs.pipeline_wrapper import pipeline_handlers
from model.crash.crash_config import CrashConfig
from libs.gitiles.change_log import ChangeLog
from libs.gitiles.change_log import FileChangeInfo
from libs.gitiles.diff import ChangeType


class ComponentClassifierTest(PredatorTestCase):
  """Tests ``ComponentClassifier`` class."""

  def setUp(self):
    super(ComponentClassifierTest, self).setUp()
    config = CrashConfig.Get().component_classifier
    components = [Component(info['component'], info['dirs'],
                            info.get('function'), info.get('team'))
                  for info in config['component_info']]
    # Only construct the classifier once, rather than making a new one every
    # time we call a method on it.
    self.classifier = ComponentClassifier(components, config['top_n'])

  def testClassifyCallStack(self):
    """Tests ``ClassifyCallStack`` method."""
    callstack = CallStack(
        0, [StackFrame(0, 'src/', 'func', 'comp1/a.cc', 'src/comp1/a.cc', [2])])
    self.assertEqual(self.classifier.ClassifyCallStack(callstack),
                     ['Comp1>Dummy'])

    callstack = CallStack(
        0, [StackFrame(0, 'dummy/', 'no_func', 'comp2/a.cc',
                       'dummy/comp2.cc', [32])])
    self.assertEqual(self.classifier.ClassifyCallStack(callstack), [])

    crash_stack = CallStack(0, frame_list=[
        StackFrame(0, 'src/', 'func', 'comp1/a.cc', 'src/comp1/a.cc', [2]),
        StackFrame(1, 'src/', 'ff', 'comp1/a.cc', 'src/comp1/a.cc', [21]),
        StackFrame(2, 'src/', 'func2', 'comp2/b.cc', 'src/comp2/b.cc', [8])])

    self.assertEqual(self.classifier.ClassifyCallStack(crash_stack),
                     ['Comp1>Dummy', 'Comp2>Dummy'])

  def testClassifySuspect(self):
    """Tests ``ClassifySuspect`` method."""
    suspect = Suspect(self.GetDummyChangeLog(), 'src/')
    suspect.changelog = suspect.changelog._replace(
        touched_files = [FileChangeInfo(ChangeType.MODIFY,
                                        'comp1/a.cc', 'comp1/b.cc')])
    self.assertEqual(self.classifier.ClassifySuspect(suspect), ['Comp1>Dummy'])

  def testClassifyEmptySuspect(self):
    """Tests ``ClassifySuspect`` returns None for empty suspect."""
    self.assertIsNone(self.classifier.ClassifySuspect(None))

  def testClassifySuspectNoMatch(self):
    """Tests ``ClassifySuspect`` returns None if there is no file match."""
    suspect = Suspect(self.GetDummyChangeLog(), 'dummy')
    suspect.changelog = suspect.changelog._replace(
        touched_files = [FileChangeInfo(ChangeType.MODIFY,
                                        'comp1.cc', 'comp1.cc')])
    self.assertEqual(self.classifier.ClassifySuspect(suspect), [])

  def testClassifySuspects(self):
    """Tests ``ClassifySuspects`` classify a list of ``Suspect``s."""
    suspect1 = Suspect(self.GetDummyChangeLog(), 'src/')
    suspect1.changelog = suspect1.changelog._replace(
        touched_files = [FileChangeInfo(ChangeType.MODIFY,
                                        'comp1/a.cc', 'comp1/b.cc')])
    suspect2 = Suspect(self.GetDummyChangeLog(), 'src/')
    suspect2.changelog = suspect2.changelog._replace(
        touched_files = [FileChangeInfo(ChangeType.MODIFY,
                                        'comp2/a.cc', 'comp2/b.cc')])

    self.assertEqual(self.classifier.ClassifySuspects([suspect1, suspect2]),
                     ['Comp1>Dummy', 'Comp2>Dummy'])
