# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""High-level wrapper of a usb device.

Wraps python-libusb1 to provide config/descriptions for a device. Does not
provide any low-level IO support over the bus.
"""

import os
import sys
import libusb1
import usb1


if sys.platform != 'linux2':
  raise NotImplementedError('This library only supported on linux systems.')


def get_all_devices():
  ctx = usb1.USBContext()
  return [USBDevice(d) for d in ctx.getDeviceList(skip_on_error=True)]


class USBDevice(object):
  def __init__(self, libusb_device):
    self._libusb_device = libusb_device

    self.port = libusb_device.getPortNumber()
    self.bus = libusb_device.getBusNumber()
    self.dev = libusb_device.getDeviceAddress()
    self.physical_port = None
    self._serial = None
    self._port_list = None

    # libusb exposes product and vendor IDs as decimal but sysfs reports
    # them as hex. Convert to hex for easy string comparison.
    self.product = hex(libusb_device.getProductID())[2:]
    self.vendor = hex(libusb_device.getVendorID())[2:]

    # libusb doesn't expose major and minor numbers, so stat the device file.
    self.major = None
    self.minor = None
    self.dev_file_path = os.path.join(
        '/dev/bus/usb', '%03d' % self.bus, '%03d' % self.dev)
    try:
      st = os.stat(self.dev_file_path)
      self.major = os.major(st.st_rdev)
      self.minor = os.minor(st.st_rdev)
    except OSError:
      pass

  def __str__(self):
    return self.serial or self.port_list

  @property
  def serial(self):
    if not self._serial:
      try:
        self._serial = self._libusb_device.getSerialNumber()
      except usb1.USBError:
        self._serial = None
    return self._serial

  @property
  def port_list(self):
    if not self._port_list:
      try:
        self._port_list = self._libusb_device.getPortNumberList()
      except usb1.USBError:
        self._port_list = None
    return self._port_list

  @property
  def interfaces(self):
    for setting in self._libusb_device.iterSettings():
      yield (setting.getClass(), setting.getSubClass(), setting.getProtocol())
