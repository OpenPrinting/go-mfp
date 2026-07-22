# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# USB-related definitions

from helpers import collection, keyword

# USB types
class DeviceDescriptor(collection): pass
class ConfigurationDescriptor(collection): pass
class Interface(collection): pass
class InterfaceDescriptor(collection): pass
class EndpointDescriptor(collection): pass

# Keywords
class IN(keyword): pass
class OUT(keyword): pass

# device is the model-settable variable that defines the
# USB device parameters.
device = None

