#!/usr/bin/env python
# Copyright (c) 2012 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

""" Compile step for Android """

from build_step import BuildStep
from utils import shell_utils
import os
import sys


class AndroidCompile(BuildStep):
  def _Run(self):
    os.environ['PATH'] = os.path.abspath(
        os.path.join(os.pardir, os.pardir, os.pardir, os.pardir, 'third_party',
                     'gsutil')) + os.pathsep + os.environ['PATH']
    os.environ['BOTO_CONFIG'] = os.path.abspath(os.path.join(
        os.pardir, os.pardir, os.pardir, os.pardir, 'site_config', '.boto'))
    os.environ['ANDROID_SDK_ROOT'] = self._args['android_sdk_root']
    os.environ['GYP_DEFINES'] = self._args['gyp_defines']
    print 'GYP_DEFINES="%s"' % os.environ['GYP_DEFINES']
    cmd = [os.path.join('platform_tools', 'android', 'bin', 'android_make'),
           self._args['target'],
           '-d', self._args['device'],
           'BUILDTYPE=%s' % self._configuration,
           ]
    cmd.extend(self._default_make_flags)
    if os.name != 'nt':
      try:
        ccache = shell_utils.Bash(['which', 'ccache'], echo=False)
        if ccache:
          cmd.append('--use-ccache')
      except Exception:
        pass
    cmd.extend(self._make_flags)
    shell_utils.Bash(cmd)


if '__main__' == __name__:
  sys.exit(BuildStep.RunBuildStep(AndroidCompile))
