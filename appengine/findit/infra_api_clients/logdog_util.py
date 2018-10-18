# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
"""Utilities for interacting with LUCI's LogDog service."""

import base64
import cStringIO
import logging
import json
import os
import sys
import time

from common import rpc_util

import google
# protobuf and GAE have package name conflict on 'google'.
# Add this to solve the conflict.
third_party = os.path.join(
    os.path.dirname(__file__), os.path.pardir, 'third_party')
sys.path.insert(0, third_party)
google.__path__.insert(0, os.path.join(third_party, 'google'))
from logdog import annotations_pb2

_LOGDOG_ENDPOINT = 'https://%s/prpc/logdog.Logs'
_LOGDOG_TAIL_ENDPOINT = '%s/Tail' % _LOGDOG_ENDPOINT
_LOGDOG_GET_ENDPOINT = '%s/Get' % _LOGDOG_ENDPOINT


def _ProcessStringForLogDog(base_string):
  """Processes base string and replaces all special characters to '_'.

  Special characters are non alphanumeric nor ':_-.'.
  Reference: https://chromium.googlesource.com/chromium/tools/build/+/refs/
      heads/master/scripts/slave/logdog_bootstrap.py#349
  """
  new_string_list = []
  for c in base_string:
    if not c.isalnum() and not c in ':_-.':
      new_string_list.append('_')
    else:
      new_string_list.append(c)
  return ''.join(new_string_list)


def _GetRawLogsFromGetEndpoint(host, data, http_client, retry_delay=5):
  """Gets raw logs from Get endpoint.

   The raw logs can be further processed to get annotations or a specific log.

   For logs of annotations, it should look like:
    [
       {
           'datagram': {
               'data': (base64 encoded data)
           }
       }
    ].

  For an actual log, it should look like:
    [
       {
           'text': {
               'lines': [
                  {
                      'value': 'line'
                  }
               ]
           }
       }
    ].
  """
  tries = 0
  error_message = ''

  # It seems possible to get empty log or log with wrong format.
  # So also retry for several times even on 200s if the log cannot be used.
  while tries < 5:
    # Retry 7 times to allow for logdog's up to 180 second propagation delay.
    # Exponential backoff starts at 1.5 seconds, reaches 96 seconds for the 7th
    # retry, for an accumulated total of 190.5 seconds of waiting time.
    # Should be enough for our purposes.
    _, response_json = rpc_util.DownloadJsonData(
        _LOGDOG_GET_ENDPOINT % host, data, http_client, max_retries=7)
    if response_json is None:
      # If response is None, it means after 7 retries, Findit still failed to
      # get response. Seems no need to keep retrying at this case.
      error_message = 'cannot get json log.'
      break
    else:
      try:
        logs = json.loads(response_json).get('logs')
        if not logs or not isinstance(logs, list):
          error_message = 'Wrong format - %s' % response_json
        else:
          return logs
      except ValueError as e:
        # For unknown reason sometimes the response_json is truncated and cannot
        # be json loaded.
        # This will also help to catch if the response_json is not serializable.
        error_message = 'Failed to load json - %s' % e.message
    tries += 1
    time.sleep(tries * retry_delay)

  # Only logs error when the log was failed to get at last.
  logging.error('Error when fetch log or annotations: %s' % error_message)
  return None


def _GetAnnotationsProtoForPath(host, project, path, http_client):
  """Gets annotations from logdog endpoint(s).

  By default sends request to Tail endpoint for annotations, if only gets a
  partial results, use Get endpoint instead.
  """
  base_error_log = 'Error when load annotations protobuf: %s'

  data = {'project': project, 'path': path}

  _, response_json = rpc_util.DownloadJsonData(
      _LOGDOG_TAIL_ENDPOINT % host, data, http_client, max_retries=7)
  if not response_json:
    return None

  # Gets data for proto. Data format as below:
  # {
  #    'logs': [
  #        {
  #            'datagram': {
  #                'data': (base64 encoded data)
  #            }
  #        }
  #     ]
  # }
  logs = json.loads(response_json).get('logs')
  if not logs or not isinstance(logs, list):
    logging.error(base_error_log % 'Wrong format - "logs"')
    return None

  partial = logs[-1].get('datagram', {}).get('partial')
  if partial:
    # Only gets partial result from Tail, use Get instead to get annotations.
    index = int(logs[-1]['streamIndex'])
    partial_index = partial['index']
    data = {'project': project, 'path': path, 'index': index - partial_index}
    logs = _GetRawLogsFromGetEndpoint(host, data, http_client)

  annotations = ''
  if not logs:
    logging.error(base_error_log % 'Wrong format - "logs"')
    return None

  sio = cStringIO.StringIO()
  for log in logs:
    annotations_b64 = log.get('datagram', {}).get('data')
    if not annotations_b64:
      sio.close()
      logging.error(base_error_log % 'Wrong format - "data"')
      return None

    sio.write(base64.b64decode(annotations_b64))
  annotations = sio.getvalue()
  sio.close()

  # Gets proto.
  try:
    step = annotations_pb2.Step()
    step.ParseFromString(annotations)
    return step
  except Exception:
    logging.error(base_error_log % 'could not get annotations.')
    return None


def _GetStreamForStep(step_name, data, log_type='stdout'):
  for substep in data.substep:
    if substep.step.name != step_name:
      continue

    if log_type.lower() == 'stdout':
      # Gets stdout_stream.
      return substep.step.stdout_stream.name

    # Gets stream for step_metadata.
    for link in substep.step.other_links:
      if link.label.lower() == log_type:
        return link.logdog_stream.name

  return None


def GetLogLocationFromBuildbucketBuild(buildbucket_build):
  """Gets the path to the logdog annotations.

  Args:
    buildbucket_build (dict): The build as retrieved by the buildbucket client.
  Returns:
    log_location of the build.
  """
  # The log location is a property on buildbot builds and a tag on swarming
  # builds.
  # First check the tags.
  for tag in buildbucket_build.get('tags', []):
    if tag.startswith('swarming_tag:log_location:logdog:'):
      return tag.split(':', 2)[2]

  # Then try to get it from properties.
  result_details = json.loads(
      buildbucket_build.get('result_details_json', '{}'))
  return result_details.get('properties', {}).get('log_location')


def _GetQueryParametersForAnnotation(log_location):
  """Gets the path to the logdog annotations.

  Args:
    log_location (str): The log location for the build.
  Returns:
    The (host, project, path) triad that identifies the location of the
    annotations proto.
  """
  host = project = path = None
  if log_location:
    # logdog://luci-logdog.appspot.com/chromium/...
    _logdog, _, host, project, path = log_location.split('/', 4)
  return host, project, path


def _GetLog(annotations, step_name, log_type, http_client):
  if not annotations:
    return None
  stream = _GetStreamForStep(step_name, annotations, log_type)
  if not stream:
    return None
  env = annotations.command.environ
  host = env['LOGDOG_COORDINATOR_HOST']
  project = env['LOGDOG_STREAM_PROJECT']
  prefix = env['LOGDOG_STREAM_PREFIX']
  if not all([host, project, prefix]):
    return None
  path = '%s/+/%s' % (prefix, stream)

  log_url = 'https://{host}/logs/{project}/{path}?format=raw'.format(
      host=host, project=project, path=path)
  status_code, log, _ = http_client.Get(log_url)

  if status_code != 200 or not log:
    logging.error('Failed to get the log from %s: status_code-%d, log-%s',
                  log_url, status_code, log)
    return None
  return log


def GetStepLogForBuild(buildbucket_build, step_name, log_type, http_client):
  log_location = GetLogLocationFromBuildbucketBuild(buildbucket_build)
  host, project, path = _GetQueryParametersForAnnotation(log_location)
  annotations = _GetAnnotationsProtoForPath(host, project, path, http_client)
  return _GetLog(annotations, step_name, log_type, http_client)


def GetStepLogLegacy(log_location, step_name, log_type, http_client):
  host, project, path = _GetQueryParametersForAnnotation(log_location)
  if not host:
    logging.error('Failed to get log_location info for logdog stream.')
    return None

  annotations = _GetAnnotationsProtoForPath(host, project, path, http_client)
  return _GetLog(annotations, step_name, log_type, http_client)
