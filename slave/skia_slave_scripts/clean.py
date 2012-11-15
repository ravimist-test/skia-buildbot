#!/usr/bin/env python
# Copyright (c) 2012 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

""" Clean step """

from utils import shell_utils
from build_step import BuildStep
import os
import sys


class Clean(BuildStep):
  def _Run(self):
    make_cmd = 'make'
    if os.name == 'nt':
      make_cmd = 'make.bat'
    shell_utils.Bash([make_cmd, 'clean'])


if '__main__' == __name__:
  sys.exit(BuildStep.RunBuildStep(Clean))