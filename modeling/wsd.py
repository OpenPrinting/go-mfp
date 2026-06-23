# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# WS-Scan definitions

from helper import collection
from dataclasses import dataclass

# WS-Scan types
class ActiveJobs(collection): pass
class ADF(collection): pass
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
class Documents(collection): pass
class Exposure(collection): pass
class ExposureSettings(collection): pass
class Film(collection): pass
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
class Platen(collection): pass
class Range(collection): pass
class Resolution(collection): pass
class Resolutions(collection): pass
class RetrieveImageRequest(collection): pass
class RetrieveImageResponse(collection): pass
class Scaling(collection): pass
class ScalingRangeSupported(collection): pass
class ScanData(collection): pass
class ScannerConfiguration(collection): pass
class ScannerDescription(collection): pass
class ScannerElemData(collection): pass
class ScannerStatus(collection): pass
class ScanRegion(collection): pass
class ScanTicket(collection): pass

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

