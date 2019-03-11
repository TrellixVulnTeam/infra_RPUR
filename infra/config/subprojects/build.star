# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Definitions of build.git CI resources."""

load('//lib/infra.star', 'infra')
load('//lib/presubmit.star', 'presubmit')
load('//lib/recipes.star', 'recipes')


REPO_URL = 'https://chromium.googlesource.com/chromium/tools/build'


infra.console_view(
    name = 'build',
    title = 'build repository console',
    repo = REPO_URL,
)

luci.cq_group(
    name = 'build cq',
    watch = cq.refset(repo = REPO_URL, refs = [r'refs/heads/.+']),
    retry_config = cq.RETRY_ALL_FAILURES,
)


# Presubmit trybots.
presubmit.builder(
    name = 'Build Presubmit',
    cq_group = 'build cq',
    repo_name = 'build',
)
# Trybot that launches a task via 'led' to verify updated recipes work.
recipes.led_recipes_tester(
    name = 'Build Recipes Tester',
    cq_group = 'build cq',
    repo_name = 'build',
)


# Recipes ecosystem.
recipes.simulation_tester(
    name = 'build-recipes-tests',
    project_under_test = 'build',
    triggered_by = 'build-gitiles-trigger',
    console_view = 'build',
)


# Recipe rolls from Build.
recipes.roll_trybots(
    upstream = 'build',
    downstream = [
        'infra',
    ],
    cq_group = 'build cq',
)


# External testers (defined in another projects) for recipe rolls.
luci.cq_tryjob_verifier(
    builder = 'external/infra-internal/try/build_limited Roll Tester (build)',
    cq_group = 'build cq',
)
luci.cq_tryjob_verifier(
    builder = 'external/infra-internal/try/release_scripts Roll Tester (build)',
    cq_group = 'build cq',
)