# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Helper methods that handle both git and svn repos."""

import infra.services.bugdroid.svn_helper as svn_helper


def GetBranch(log_entry, full=False):
  if log_entry.scm == 'svn':
    return svn_helper.get_branch(log_entry, full=full)
  elif log_entry.scm == 'git':
    # Strip off the prefix (e.g. 'refs/heads/') to make the returned value more
    # like what svn returns, since this is used in things like merge-merged
    # issue labels, which are, in turn, probably used in places that expect
    # merge labels to just have a "-#####" suffix.
    if log_entry.branch:
      if full:
        return log_entry.branch
      else:
        return log_entry.branch.split('/')[-1]
  return None