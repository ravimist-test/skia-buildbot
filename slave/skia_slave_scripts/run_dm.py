#!/usr/bin/env python
# Copyright (c) 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Run the Skia DM executable. """

from build_step import BuildStep
import sys

class RunDM(BuildStep):
  def _AnyMatch(self, *args):
    return any(arg in self._builder_name for arg in args)

  def _AllMatch(self, *args):
    return all(arg in self._builder_name for arg in args)

  def _Run(self):
    args = [
      '--verbose',
      '--resourcePath', self._device_dirs.ResourceDir(),
    ]

    run_dm = True
    match  = []

    # Disable failing tests.
    if self._AnyMatch('ChromeOS', 'DirectWrite'):
      match.append('~bitmapscroll')

    if self._AnyMatch('Win7', 'Android'):
      match.append('~blurroundrect')

    if self._AnyMatch('Android'):
      match.append('~giantbitmap')

    if self._AnyMatch('Tegra'):
      match.append('~downsamplebitmap_text')

    # Disable crashing tests.
    if self._AllMatch('10.6', 'Debug'):
      # Not sure what's failing exactly, so disable DM entirely.
      run_dm = False

    if match:
      args.append('--match')
      args.extend(match)

    if run_dm:
      self._flavor_utils.RunFlavoredCmd('dm', args)

if '__main__' == __name__:
  sys.exit(BuildStep.RunBuildStep(RunDM))
