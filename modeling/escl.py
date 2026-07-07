# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# eSCL-related definitions

from uuid import UUID
from typing import TypedDict
from helpers import collection, keyword

# eSCL types
class Adf(collection): pass
class Camera(collection, keyword): pass
class DiscreteResolution(collection): pass
class InputSourceCaps(collection): pass
class JobInfo(collection): pass
class Justification(collection): pass
class Platen(collection, keyword): pass
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

# Keywords
class AbortedBySystem (keyword): pass
class Aborted (keyword): pass
class AccountAuthorizationFailed (keyword): pass
class AccountClosed (keyword): pass
class AccountInfoNeeded (keyword): pass
class AccountLimitReached (keyword): pass
class Auto (keyword): pass
class BlackAndWhite1 (keyword): pass
class Blue (keyword): pass
class BottomEdge (keyword): pass
class Bottom (keyword): pass
class BusinessCard (keyword): pass
class Canceled (keyword): pass
class Center (keyword): pass
class Completed (keyword): pass
class CompressionError (keyword): pass
class ConflictingAttributes (keyword): pass
class ConnectedToDestination (keyword): pass
class ConnectingToDestination (keyword): pass
class DestinationUriFailed (keyword): pass
class DetectPaperLoaded (keyword): pass
class DigitalSignatureDidNotVerify (keyword): pass
class DigitalSignatureTypeNotSupported (keyword): pass
class DocumentAccessError (keyword): pass
class DocumentFormatError (keyword): pass
class Document (keyword): pass
class DocumentPasswordError (keyword): pass
class DocumentPermissionError (keyword): pass
class DocumentSecurityError (keyword): pass
class DocumentUnprintableError (keyword): pass
class Down (keyword): pass
class Duplex (keyword): pass
class ErrorsDetected (keyword): pass
class Feeder (keyword): pass
class GrayCcdEmulated (keyword): pass
class GrayCcd (keyword): pass
class Grayscale16 (keyword): pass
class Grayscale8 (keyword): pass
class Green (keyword): pass
class Halftone (keyword): pass
class Idle (keyword): pass
class JobCanceledAtDevice (keyword): pass
class JobCanceledByOperator (keyword): pass
class JobCanceledByUser (keyword): pass
class JobCompletedSuccessfully (keyword): pass
class JobCompletedWithErrors (keyword): pass
class JobCompletedWithWarnings (keyword): pass
class JobDataInsufficient (keyword): pass
class JobDelayOutputUntilSpecified (keyword): pass
class JobDigitalSignatureWait (keyword): pass
class JobFetchable (keyword): pass
class JobHeldByService (keyword): pass
class JobHeldForReview (keyword): pass
class JobHoldUntilSpecified (keyword): pass
class JobIncoming (keyword): pass
class JobInterpreting (keyword): pass
class JobOutgoing (keyword): pass
class JobPasswordWait (keyword): pass
class JobPrintedSuccessfully (keyword): pass
class JobPrintedWithErrors (keyword): pass
class JobPrintedWithWarnings (keyword): pass
class JobPrinting (keyword): pass
class JobQueuedForMarker (keyword): pass
class JobQueued (keyword): pass
class JobReleaseWait (keyword): pass
class JobRestartable (keyword): pass
class JobResuming (keyword): pass
class JobSavedSuccessfully (keyword): pass
class JobSavedWithErrors (keyword): pass
class JobSavedWithWarnings (keyword): pass
class JobSaving (keyword): pass
class JobScanningAndTransferring (keyword): pass
class JobScanning (keyword): pass
class JobSpooling (keyword): pass
class JobStreaming (keyword): pass
class JobSuspendedByOperator (keyword): pass
class JobSuspendedBySystem (keyword): pass
class JobSuspendedByUser (keyword): pass
class JobSuspended (keyword): pass
class JobSuspending (keyword): pass
class JobTransferring (keyword): pass
class JobTransforming (keyword): pass
class LeftEdge (keyword): pass
class Left (keyword): pass
class LineArt (keyword): pass
class LongEdgeFeed (keyword): pass
class Magazine (keyword): pass
class NTSC (keyword): pass
class Object (keyword): pass
class PendingHeld (keyword): pass
class Pending (keyword): pass
class Photo (keyword): pass
class Preview (keyword): pass
class PrinterStopped (keyword): pass
class PrinterStoppedPartly (keyword): pass
class Processing (keyword): pass
class ProcessingToStopPoint (keyword): pass
class QueuedInDevice (keyword): pass
class Red (keyword): pass
class ResourcesAreNotReady (keyword): pass
class ResourcesAreNotSupported (keyword): pass
class RGB24 (keyword): pass
class RGB48 (keyword): pass
class RightEdge (keyword): pass
class Right (keyword): pass
class ScannerAdfDuplexPageTooLong (keyword): pass
class ScannerAdfDuplexPageTooShort (keyword): pass
class ScannerAdfEmpty (keyword): pass
class ScannerAdfHatchOpen (keyword): pass
class ScannerAdfInputTrayFailed (keyword): pass
class ScannerAdfInputTrayOverloaded (keyword): pass
class ScannerAdfJam (keyword): pass
class ScannerAdfLoaded (keyword): pass
class ScannerAdfMispick (keyword): pass
class ScannerAdfMultipickDetected (keyword): pass
class ScannerAdfProcessing (keyword): pass
class SelectSinglePage (keyword): pass
class ServiceOffLine (keyword): pass
class ShortEdgeFeed (keyword): pass
class sRGB (keyword): pass
class Stopped (keyword): pass
class SubmissionInterrupted (keyword): pass
class Testing (keyword): pass
class TextAndGraphic (keyword): pass
class TextAndPhoto (keyword): pass
class Text (keyword): pass
class ThreeHundredthsOfInches (keyword): pass
class Threshold (keyword): pass
class TopEdge (keyword): pass
class Top (keyword): pass
class UnsupportedAttributesOrValues (keyword): pass
class UnsupportedCompression (keyword): pass
class UnsupportedDocumentFormat (keyword): pass
class WaitingForUserAction (keyword): pass
class WarningsDetected (keyword): pass

class ImageFilter(TypedDict):
    OutputFormat: str
    XResolution: int
    YResolution: int
    ColorMode: str

# caps is the model-settable variable that defines the
# eSCL scanner capabilities
caps = None

