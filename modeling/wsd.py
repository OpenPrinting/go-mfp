# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# WS-Scan definitions

from helpers import collection, keyword, iskeyword
from dataclasses import dataclass

# WS-Scan types
class ActiveJobs(collection): pass
class ADF(collection, keyword): pass
class ADFSide(collection): pass
class CancelJobRequest(collection): pass
class CancelJobResponse(collection): pass
class ConditionHistoryEntry(collection): pass
class CreateScanJobRequest(collection): pass
class CreateScanJobResponse(collection): pass
class DeviceCondition(collection): pass
class DeviceSettings(collection): pass
class Dimensions(collection): pass
class Document(collection): pass
class DocumentDescription(collection): pass
class DocumentParameters(collection): pass
class Documents(collection, keyword): pass
class Exposure(collection): pass
class ExposureSettings(collection): pass
class Film(collection, keyword): pass
class GetActiveJobsRequest(collection): pass
class GetActiveJobsResponse(collection): pass
class GetJobElementsRequest(collection): pass
class GetJobElementsResponse(collection): pass
class GetJobHistoryRequest(collection): pass
class GetJobHistoryResponse(collection): pass
class GetScannerElementsRequest(collection): pass
class GetScannerElementsResponse(collection): pass
class ImageInformation(collection): pass
class InputMediaSize(collection): pass
class InputSize(collection): pass
class Job(collection): pass
class JobDescription(collection): pass
class JobElemData(collection): pass
class JobStatus(collection): pass
class JobSummary(collection): pass
class MediaSide(collection): pass
class MediaSideImageInfo(collection): pass
class MediaSides(collection): pass
class Platen(collection, keyword): pass
class Range(collection): pass
class Resolution(collection): pass
class Resolutions(collection): pass
class RetrieveImageRequest(collection): pass
class RetrieveImageResponse(collection): pass
class Scaling(collection): pass
class ScalingRangeSupported(collection): pass
class ScanData(collection): pass
class ScannerConfiguration(collection, keyword): pass
class ScannerDescription(collection, keyword): pass
class ScannerElemData(collection): pass
class ScannerStatus(collection, keyword): pass
class ScanRegion(collection): pass
class ScanTicket(collection): pass

# Keywords
class Aborted (keyword): pass
class ADFDuplex (keyword): pass
class AttentionRequired (keyword): pass
class Auto (keyword): pass
class BlackAndWhite1 (keyword): pass
class BlackandWhiteNegativeFilm (keyword): pass
class Calibrating (keyword): pass
class Canceled (keyword): pass
class ColorNegativeFilm (keyword): pass
class ColorSlideFilm (keyword): pass
class Completed (keyword): pass
class CoverOpen (keyword): pass
class Creating (keyword): pass
class Critical (keyword): pass
class DefaultScanTicket (keyword): pass
class DocumentFormatError (keyword): pass
class Grayscale16 (keyword): pass
class Grayscale4 (keyword): pass
class Grayscale8 (keyword): pass
class Halftone (keyword): pass
class Held (keyword): pass
class Idle (keyword): pass
class ImageTransferError (keyword): pass
class Informational (keyword): pass
class InputTrayEmpty (keyword): pass
class InterlockOpen (keyword): pass
class InternalStorageFull (keyword): pass
class InvalidScanTicket (keyword): pass
class JobCanceledAtDevice (keyword): pass
class JobCompletedWithErrors (keyword): pass
class JobCompletedWithWarnings (keyword): pass
class JobScanningAndTransferring (keyword): pass
class JobScanning (keyword): pass
class JobTimedOut (keyword): pass
class JobTransferring (keyword): pass
class LampError (keyword): pass
class LampWarming (keyword): pass
class MediaJam (keyword): pass
class MediaPath (keyword): pass
class Mixed (keyword): pass
class MultipleFeedError (keyword): pass
class NotApplicable (keyword): pass
class Paused (keyword): pass
class Pending (keyword): pass
class Photo (keyword): pass
class Processing (keyword): pass
class RGB24 (keyword): pass
class RGB48 (keyword): pass
class RGBa32 (keyword): pass
class RGBa64 (keyword): pass
class ScannerStopped (keyword): pass
class Started (keyword): pass
class Stopped (keyword): pass
class Terminating (keyword): pass
class Text (keyword): pass
class VendorSection (keyword): pass
class Warning (keyword): pass

# Boolean represents the WSD boolean value.
#
# This is considered True if equal to "1" or "true" (case-insensitive),
# false otherwise.
class Boolean(str):
    def __bool__(self):
        lwr = self.lower()
        return lwr == "1" or lwr == "true"

# WithOptions wraps arbitrary class into dynamically created
# class with additional Boolean options: MustHonor, Override
# and UsedDefault
def WithOptions(value, *, MustHonor: Boolean = None, Override: Boolean = None, UsedDefault: Boolean = None):
    # Save value's representation
    base_repr = repr(value)

    # WithOptions doesn't work well with wrapping classes,
    # and keywords are classes. So for keyword, replace
    # value with its string representation.
    if iskeyword(value):
        value = str(value)

    base_type = type(value)

    # Wrap value into DynamicWithOptions
    class DynamicWithOptions(base_type):
        def __init__(self, val, *args, **kwargs):
            self.MustHonor = MustHonor
            self.Override = Override
            self.UsedDefault = UsedDefault

        def __repr__(self):
            # Gather active options
            options = {
                'MustHonor': self.MustHonor,
                'Override': self.Override,
                'UsedDefault': self.UsedDefault
            }


            active_options = {k: v for k, v in options.items() if v is not None}

            # Format base value
            #base_repr = super().__repr__()

            if not active_options:
                return base_repr

            # Format WithOptions call
            opts_str = ", ".join(f"{k}={repr(v)}" for k, v in active_options.items())
            return f"WithOptions({base_repr}, {opts_str})"

    DynamicWithOptions.__name__ = 'WithOptions'
    DynamicWithOptions.__qualname__ = 'WithOptions'

    return DynamicWithOptions(value)

# WithLang represents a string with optional language identifier.
@dataclass(repr=False, init=False)
class WithLang(str):
    lang: str = None

    def __new__(cls, value: str, lang: str = None):
        v = super().__new__(cls, value)
        v.lang = lang
        return v

    def __repr__(self) -> str:
        if self.lang is None:
            return f"'{self}'"
        return f"wsd.WithLang('{self}', lang='{self.lang}')"


# caps is the model-settable variable that defines the
# WS-Scan scanner capabilities
caps = None

