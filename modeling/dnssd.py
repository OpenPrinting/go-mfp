# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# DNS-SD-related definitions

from helpers import collection

class Device(collection): pass
class Service(collection): pass

# device is the model-settable variable that defines the
# printer's DNS-SD parameters.
device = None

