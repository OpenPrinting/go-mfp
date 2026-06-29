# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# eSCL-related definitions

from uuid import UUID
from typing import TypedDict
from helpers import collection

# eSCL types
class Adf(collection): pass
class Camera(collection): pass
class DiscreteResolution(collection): pass
class InputSourceCaps(collection): pass
class JobInfo(collection): pass
class Platen(collection): pass
class Range(collection): pass
class Region(collection): pass
class ResolutionRange(collection): pass
class ScanBufferInfo(collection): pass
class ScanImageInfo(collection): pass
class ScannerCapabilities(collection): pass
class ScannerStatus(collection): pass
class ScanRegion(collection): pass
class ScanSettings(collection): pass
class SettingProfile(collection): pass
class SupportedResolutions(collection): pass

class ImageFilter(TypedDict):
    OutputFormat: str
    XResolution: int
    YResolution: int
    ColorMode: str

# caps is the model-settable variable that defines the
# eSCL scanner capabilities
caps = None

