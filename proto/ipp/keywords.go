// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP keywords

package ipp

import (
	"reflect"
	"strings"
)

// KwCompression represents standard keyword values for compression.
//
// See RFC8011, 5.4.32.
type KwCompression string

const (
	// KwCompressionNone is no compression
	KwCompressionNone KwCompression = "none"

	// KwCompressionDeflate is RFC1951 ZIP inflate/deflate
	KwCompressionDeflate KwCompression = "deflate"

	// KwCompressionGzip is RFC1952 GNU zip
	KwCompressionGzip KwCompression = "gzip"

	// KwCompressionCompress is RFC1977 UNIX compression
	KwCompressionCompress KwCompression = "compress"
)

// KwJobDelayOutputUntil represents standard keyword values for
// "job-delay-output-until" attribute.
//
// See PWG5100.7, 6.8.4.
type KwJobDelayOutputUntil string

const (
	// KwJobDelayOutputUntilDayTime means delay output until daylight,
	// typically 6am to 6pm.
	KwJobDelayOutputUntilDayTime KwJobDelayOutputUntil = "day-time"

	// KwJobDelayOutputUntilEvening means delay output until the evening,
	// typically from 6pm to 12am.
	KwJobDelayOutputUntilEvening KwJobDelayOutputUntil = "evening"

	// KwJobDelayOutputUntilIndefinite means delay output indefinitely;
	// the time period can be changed using the Set- Job-Attributes
	// operation
	KwJobDelayOutputUntilIndefinite KwJobDelayOutputUntil = "indefinite"

	// KwJobDelayOutputUntilNight means delay output until the night,
	// typically from 12am to 6am.
	KwJobDelayOutputUntilNight KwJobDelayOutputUntil = "night"

	// KwJobDelayOutputUntilNoDelayOutput means do not delay the output.
	KwJobDelayOutputUntilNoDelayOutput KwJobDelayOutputUntil = "no-delay-output"

	// KwJobDelayOutputUntilSecondShift means delay output until the
	// second work shift, typically from 4pm to 12am.
	KwJobDelayOutputUntilSecondShift KwJobDelayOutputUntil = "second-shift"

	// KwJobDelayOutputUntilThirdShift means delay output until the
	// third work shift, typically from 12am to 8am.
	KwJobDelayOutputUntilThirdShift KwJobDelayOutputUntil = "third-shift"

	// KwJobDelayOutputUntilWeekend means delay output until the weekend.
	KwJobDelayOutputUntilWeekend KwJobDelayOutputUntil = "weekend"
)

// KwJobHoldUntil represents standard keyword values for
// "job-hold-until" attribute.
//
// See PWG5100.7, 5.2.2.
type KwJobHoldUntil string

const (
	// KwJobHoldUntilNoHold means do not delay the output.
	KwJobHoldUntilNoHold KwJobHoldUntil = "no-hold"

	// KwJobHoldUntilIndefinite means delay output indefinitely;
	// the time period can be changed using the Set- Job-Attributes
	// operation
	KwJobHoldUntilIndefinite KwJobHoldUntil = "indefinite"

	// KwJobHoldUntilDayTime means delay output until daylight,
	// typically 6am to 6pm.
	KwJobHoldUntilDayTime KwJobHoldUntil = "day-time"

	// KwJobHoldUntilEvening means delay output until the evening,
	// typically from 6pm to 12am.
	KwJobHoldUntilEvening KwJobHoldUntil = "evening"

	// KwJobHoldUntilNight means delay output until the night,
	// typically from 12am to 6am.
	KwJobHoldUntilNight KwJobHoldUntil = "night"

	// KwJobHoldUntilWeekend means delay output until the weekend.
	KwJobHoldUntilWeekend KwJobHoldUntil = "weekend"

	// KwJobHoldUntilSecondShift means delay output until the
	// second work shift, typically from 4pm to 12am.
	KwJobHoldUntilSecondShift KwJobHoldUntil = "second-shift"

	// KwJobHoldUntilThirdShift means delay output until the
	// third work shift, typically from 12am to 8am.
	KwJobHoldUntilThirdShift KwJobHoldUntil = "third-shift"
)

// KwJobSheets represents standard keyword values for
// "job-sheets" attribute.
//
// See: RFC8011, 5.2.3., PWG5100.7, 8.1.
type KwJobSheets string

const (
	// ----- RFC8011 values -----

	// KwJobSheetsNone means no Job sheet is printed
	KwJobSheetsNone KwJobSheets = "none"

	// KwJobSheetsStandard means one or more site-specific
	// standard Job sheets are printed
	KwJobSheetsStandard KwJobSheets = "standard"

	// ----- PWG5100.7 values -----

	// KwJobSheetsJobStartSheet means that a ob Sheet is printed to
	// indicate the start of the Job.
	KwJobSheetsJobStartSheet KwJobSheets = "job-start-sheet "

	// KwJobSheetsJobEndSheet means that a ob Sheet is printed to
	// indicate the end of the Job.
	KwJobSheetsJobEndSheet KwJobSheets = "job-end-sheet"

	// KwJobSheetsJobBothSheets instructs Printer to print
	// Job Sheets o indicate the start and end of all the
	// output associated with the Job.
	KwJobSheetsJobBothSheets KwJobSheets = "job-both-sheets"

	// KwJobSheetsFirstPrintStreamPage instructs Printer to print the
	// first input Page in the Document Data as the Job Sheet.
	// The Printer's standard Job Sheet is suppressed.
	KwJobSheetsFirstPrintStreamPage KwJobSheets = "first-print-stream-page"
)

// KwJobSpooling represents standard keyword values for
// "job-spooling-supported" attribute.
//
// See PWG5100.7, 6.9.31.
type KwJobSpooling string

const (
	// KwJobSpoolingAutomatic means that it is up to the
	// Printer when to process the Document data.
	KwJobSpoolingAutomatic KwJobSpooling = "automatic"

	// KwJobSpoolingSpool means that the Document data is
	// processed after it has been spooled (stored).
	KwJobSpoolingSpool KwJobSpooling = "spool"

	// KwJobSpoolingStream means that the Document data is
	// processed as it is received.
	KwJobSpoolingStream KwJobSpooling = "stream"
)

// KwJobStateReasons represents standard keyword values for
// "job-state-reasons" attribute.
//
// RFC8011: 5.3.8.
type KwJobStateReasons string

// Standard values for KwJobStateReasons attribute.
const (
	// rfc3998: 9.1.
	KwJobStateReasonsJobSuspended KwJobStateReasons = "job-suspended"

	// RFC8011: 5.3.8.
	KwJobStateReasonsNone KwJobStateReasons = "none"

	KwJobStateReasonsAbortedBySystem           KwJobStateReasons = "aborted-by-system"
	KwJobStateReasonsCompressionError          KwJobStateReasons = "compression-error"
	KwJobStateReasonsDocumentAccessError       KwJobStateReasons = "document-access-error"
	KwJobStateReasonsDocumentFormatError       KwJobStateReasons = "document-format-error"
	KwJobStateReasonsJobCanceledAtDevice       KwJobStateReasons = "job-canceled-at-device"
	KwJobStateReasonsJobCanceledByOperator     KwJobStateReasons = "job-canceled-by-operator"
	KwJobStateReasonsJobCanceledByUser         KwJobStateReasons = "job-canceled-by-user"
	KwJobStateReasonsJobCompletedSuccessfully  KwJobStateReasons = "job-completed-successfully"
	KwJobStateReasonsJobCompletedWithErrors    KwJobStateReasons = "job-completed-with-errors"
	KwJobStateReasonsJobCompletedWithWarnings  KwJobStateReasons = "job-completed-with-warnings"
	KwJobStateReasonsJobDataInsufficient       KwJobStateReasons = "job-data-insufficient"
	KwJobStateReasonsJobHoldUntilSpecified     KwJobStateReasons = "job-hold-until-specified"
	KwJobStateReasonsJobIncoming               KwJobStateReasons = "job-incoming"
	KwJobStateReasonsJobInterpreting           KwJobStateReasons = "job-interpreting"
	KwJobStateReasonsJobOutgoing               KwJobStateReasons = "job-outgoing"
	KwJobStateReasonsJobPrinting               KwJobStateReasons = "job-printing"
	KwJobStateReasonsJobQueuedForMarker        KwJobStateReasons = "job-queued-for-marker"
	KwJobStateReasonsJobQueued                 KwJobStateReasons = "job-queued"
	KwJobStateReasonsJobRestartable            KwJobStateReasons = "job-restartable"
	KwJobStateReasonsJobTransforming           KwJobStateReasons = "job-transforming"
	KwJobStateReasonsPrinterStoppedPartly      KwJobStateReasons = "printer-stopped-partly"
	KwJobStateReasonsPrinterStopped            KwJobStateReasons = "printer-stopped"
	KwJobStateReasonsProcessingToStopPoint     KwJobStateReasons = "processing-to-stop-point"
	KwJobStateReasonsQueuedInDevice            KwJobStateReasons = "queued-in-device"
	KwJobStateReasonsResourcesAreNotReady      KwJobStateReasons = "resources-are-not-ready"
	KwJobStateReasonsServiceOffLine            KwJobStateReasons = "service-off-line"
	KwJobStateReasonsSubmissionInterrupted     KwJobStateReasons = "submission-interrupted"
	KwJobStateReasonsUnsupportedCompression    KwJobStateReasons = "unsupported-compression"
	KwJobStateReasonsUnsupportedDocumentFormat KwJobStateReasons = "unsupported-document-format"

	// PWG5100.3: 6.1
	KwJobStateReasonsResourcesAreNotSupported KwJobStateReasons = "resources-are-not-supported"

	// PWG5100.7: 8.2, 11.1.
	KwJobStateReasonsDigitalSignatureDidNotVerify     KwJobStateReasons = "digital-signature-did-not-verify"
	KwJobStateReasonsDigitalSignatureTypeNotSupported KwJobStateReasons = "digital-signature-type-not-supported"
	KwJobStateReasonsErrorsDetected                   KwJobStateReasons = "errors-detected"
	KwJobStateReasonsJobDelayOutputUntilSpecified     KwJobStateReasons = "job-delay-output-until-specified"
	KwJobStateReasonsJobDigitalSignatureWait          KwJobStateReasons = "job-digital-signature-wait"
	KwJobStateReasonsJobSpooling                      KwJobStateReasons = "job-spooling"
	KwJobStateReasonsJobStreaming                     KwJobStateReasons = "job-streaming"
	KwJobStateReasonsWarningsDetected                 KwJobStateReasons = "warnings-detected"

	// PWG5100.11: 11.3.
	KwJobStateReasonsjobPasswordWait        KwJobStateReasons = "job-password-wait"
	KwJobStateReasonsjobPrintedSuccessfully KwJobStateReasons = "job-printed-successfully"
	KwJobStateReasonsjobPrintedWithErrors   KwJobStateReasons = "job-printed-with-errors"
	KwJobStateReasonsjobPrintedWithWarnings KwJobStateReasons = "job-printed-with-warnings"
	KwJobStateReasonsjobResuming            KwJobStateReasons = "job-resuming"
	KwJobStateReasonsjobSavedSuccessfully   KwJobStateReasons = "job-saved-successfully"
	KwJobStateReasonsjobSavedWithErrors     KwJobStateReasons = "job-saved-with-errors"
	KwJobStateReasonsjobSavedWithWarnings   KwJobStateReasons = "job-saved-with-warnings"
	KwJobStateReasonsjobSaving              KwJobStateReasons = "job-saving"
	KwJobStateReasonsjobSuspendedByOperator KwJobStateReasons = "job-suspended-by-operator"
	KwJobStateReasonsjobSuspendedBySystem   KwJobStateReasons = "job-suspended-by-system"
	KwJobStateReasonsjobSuspendedByUser     KwJobStateReasons = "job-suspended-by-user"
	KwJobStateReasonsjobSuspending          KwJobStateReasons = "job-suspending"

	// PWG5100.13: 9.1.
	KwJobStateReasonsDocumentPasswordError    KwJobStateReasons = "document-password-error"
	KwJobStateReasonsDocumentPermissionError  KwJobStateReasons = "document-permission-error"
	KwJobStateReasonsDocumentSecurityError    KwJobStateReasons = "document-security-error"
	KwJobStateReasonsDocumentUnprintableError KwJobStateReasons = "document-unprintable-error"

	// PWG5100.15: 8.2
	KwJobStateReasonsConnectedToDestination    KwJobStateReasons = "connected-to-destination"
	KwJobStateReasonsConnectingToDestination   KwJobStateReasons = "connecting-to-destination"
	KwJobStateReasonsDestinationURIFailed      KwJobStateReasons = "destination-uri-failed"
	KwJobStateReasonsFaxModemCarrierLost       KwJobStateReasons = "fax-modem-carrier-lost"
	KwJobStateReasonsFaxModemEquipmentFailure  KwJobStateReasons = "fax-modem-equipment-failure"
	KwJobStateReasonsFaxModemInactivityTimeout KwJobStateReasons = "fax-modem-inactivity-timeout"
	KwJobStateReasonsFaxModemLineBusy          KwJobStateReasons = "fax-modem-line-busy"
	KwJobStateReasonsFaxModemNoAnswer          KwJobStateReasons = "fax-modem-no-answer"
	KwJobStateReasonsFaxModemNoDialTone        KwJobStateReasons = "fax-modem-no-dial-tone"
	KwJobStateReasonsFaxModemProtocolError     KwJobStateReasons = "fax-modem-protocol-error"
	KwJobStateReasonsFaxModemTrainingFailure   KwJobStateReasons = "fax-modem-training-failure"
	KwJobStateReasonsFaxModemVoiceDetected     KwJobStateReasons = "fax-modem-voice-detected"
	KwJobStateReasonsJobTransferring           KwJobStateReasons = "job-transferring"

	// PWG5100.16: 8.1
	KwJobStateReasonsAccountAuthorizationFailed    KwJobStateReasons = "account-authorization-failed"
	KwJobStateReasonsAccountClosed                 KwJobStateReasons = "account-closed"
	KwJobStateReasonsAccountInfoNeeded             KwJobStateReasons = "account-info-needed"
	KwJobStateReasonsAccountLimitReached           KwJobStateReasons = "account-limit-reached"
	KwJobStateReasonsConflictingAttributes         KwJobStateReasons = "conflicting-attributes"
	KwJobStateReasonsJobHeldForReview              KwJobStateReasons = "job-held-for-review"
	KwJobStateReasonsJobReleaseWait                KwJobStateReasons = "job-release-wait"
	KwJobStateReasonsUnsupportedAttributesOrValues KwJobStateReasons = "unsupported-attributes-or-values"

	// PWG5100.17: 9.3.
	KwJobStateReasonsWaitingForUserAction KwJobStateReasons = "waiting-for-user-action"

	// PWG5100.18: 9.4.
	KwJobStateReasonsJobFetchable KwJobStateReasons = "job-fetchable"
)

// KwMediaBackCoating represents standard keyword values for
// "media-back-coating" attribute.
//
// PWG5100.7: 6.3.1.2
type KwMediaBackCoating string

const (
	// KwMediaBackCoatingNone means the media does not have any coating.
	KwMediaBackCoatingNone KwMediaBackCoating = "none"

	// KwMediaBackCoatingGlossy means the media has a "glossy" coating.
	KwMediaBackCoatingGlossy KwMediaBackCoating = "glossy"

	// KwMediaBackCoatingHighGloss means the media has a "high-gloss"
	// coating.
	KwMediaBackCoatingHighGloss KwMediaBackCoating = "high-gloss"

	// KwMediaBackCoatingSemiGloss means the media has a "semi-gloss"
	// coating.
	KwMediaBackCoatingSemiGloss KwMediaBackCoating = "semi-gloss"

	// KwMediaBackCoatingSatin means the media has a "satin" coating.
	KwMediaBackCoatingSatin KwMediaBackCoating = "satin"

	// KwMediaBackCoatingMatte means the media has a "matte" coating.
	KwMediaBackCoatingMatte KwMediaBackCoating = "matte"
)

// KwMultipleDocumentHandling represents standard keyword values for
// "multiple-document-handling" attribute.
//
// See RFC8011, 5.2.4.
type KwMultipleDocumentHandling string

// KwMultipleDocumentHandling constants indicate, how Printer handles
// Job with multiple documents.
//
// Imagine we have a Job with documents "a" and "b". Also, we may
// request a single copy or multiple copies. If multiple copies of,
// say, document "a" are requested,  we will denote it as a(*).
//
// So, depending of the "multiple-document-handling" parameter,
// Job with multiple documents "a" and "b" will be printed
// as follows (square brakes means finishing process):
//
//	"single-document"
//	Documents are concatenated,  and "b" may start at the last page of "a",
//	but each copy of a+b concatenation starts at its own page.
//	  Sungle copy:       [a+b]
//	  Multiple copies:   [a+b] [a+b] ... [a+b]
//
//	"single-document-new-sheet"
//	Like "single-document", but "b" starts at its own page:
//	  Single copy:       [a b]
//	  Multiple copies:   [a b] [a b] ... [a b]
//
//	"separate-documents-uncollated-copies"
//	Each document handled separately. First printed all copies
//	of "a", then all copies of "b":
//	  Single copy:       [a] [b]
//	  Multiple copies:   [a] [a] ... [a]  [b] [b] ... [b]
//
//	"separate-documents-collated-copies"
//	Like "separate-documents-uncollated-copies", but in a case
//	of multiple copies, ordering is different:
//	  Single copy:       [a] [b]
//	  Multiple copies:   [a] [b] [a] [b] ... [a] [b]
const (
	KwMultipleDocumentHandlingSingleDocument         KwMultipleDocumentHandling = "single-document"
	KwMultipleDocumentHandlingSingleDocumentNewSheet KwMultipleDocumentHandling = "single-document-new-sheet"

	KwMultipleDocumentHandlingSeparateDocumentsUncollatedCopies KwMultipleDocumentHandling = "separate-documents-uncollated-copies"
	KwMultipleDocumentHandlingSeparateDocumentsCollatedCopies   KwMultipleDocumentHandling = "separate-documents-collated-copies"
)

// KwPdlOverride represents standard keyword values for
// "pdl-override-supported" attribute.
//
// See RFC8011, 5.4.28.
type KwPdlOverride string

const (
	// KwPdlOverrideAattempted indicates that Printer attempts to
	// make the IPP attribute values take precedence over embedded
	// instructions in the Document data.
	KwPdlOverrideAattempted KwPdlOverride = "attempted"

	// KwPdlOverrideNotAttempted indicates that the Printer makes no
	// attempt to make the IPP attribute values take precedence over
	// embedded instructions in the Document data.
	KwPdlOverrideNotAttempted KwPdlOverride = "not-attempted"
)

// KwPrinterStateReasons represents standard keyword values for
// "printer-state-reasons" attribute.
//
// See RFC8011, 5.4.12.
type KwPrinterStateReasons string

// Standard values for KwPrinterStateReasons
const (
	// The following are standard keyword values of "printer-state-reasons"
	// attribute:
	KwPrinterStateNone                           KwPrinterStateReasons = "none"
	KwPrinterStateOther                          KwPrinterStateReasons = "other"
	KwPrinterStateConnectingToDevice             KwPrinterStateReasons = "connecting-to-device"
	KwPrinterStateCoverOpen                      KwPrinterStateReasons = "cover-open"
	KwPrinterStateDeveloperEmpty                 KwPrinterStateReasons = "developer-empty"
	KwPrinterStateDeveloperLow                   KwPrinterStateReasons = "developer-low"
	KwPrinterStateDoorOpen                       KwPrinterStateReasons = "door-open"
	KwPrinterStateFuserOverTemp                  KwPrinterStateReasons = "fuser-over-temp"
	KwPrinterStateFuserUnderTemp                 KwPrinterStateReasons = "fuser-under-temp"
	KwPrinterStateInputTrayMissing               KwPrinterStateReasons = "input-tray-missing"
	KwPrinterStateInterlockOpen                  KwPrinterStateReasons = "interlock-open"
	KwPrinterStateInterpreterResourceUnavailable KwPrinterStateReasons = "interpreter-resource-unavailable"
	KwPrinterStateMarkerSupplyEmpty              KwPrinterStateReasons = "marker-supply-empty"
	KwPrinterStateMarkerSupplyLow                KwPrinterStateReasons = "marker-supply-low"
	KwPrinterStateMarkerWasteAlmostFull          KwPrinterStateReasons = "marker-waste-almost-full"
	KwPrinterStateMarkerWasteFull                KwPrinterStateReasons = "marker-waste-full"
	KwPrinterStateMediaEmpty                     KwPrinterStateReasons = "media-empty"
	KwPrinterStateMediaJam                       KwPrinterStateReasons = "media-jam"
	KwPrinterStateMediaLow                       KwPrinterStateReasons = "media-low"
	KwPrinterStateMediaNeeded                    KwPrinterStateReasons = "media-needed"
	KwPrinterStateMovingToPaused                 KwPrinterStateReasons = "moving-to-paused"
	KwPrinterStateOpcLifeOver                    KwPrinterStateReasons = "opc-life-over"
	KwPrinterStateOpcNearEol                     KwPrinterStateReasons = "opc-near-eol"
	KwPrinterStateOutputAreaAlmostFull           KwPrinterStateReasons = "output-area-almost-full"
	KwPrinterStateOutputAreaFull                 KwPrinterStateReasons = "output-area-full"
	KwPrinterStateOutputTrayMissing              KwPrinterStateReasons = "output-tray-missing"
	KwPrinterStatePaused                         KwPrinterStateReasons = "paused"
	KwPrinterStateShutdown                       KwPrinterStateReasons = "shutdown"
	KwPrinterStateSpoolAreaFull                  KwPrinterStateReasons = "spool-area-full"
	KwPrinterStateStoppedPartly                  KwPrinterStateReasons = "stopped-partly"
	KwPrinterStateStopping                       KwPrinterStateReasons = "stopping"
	KwPrinterStateTimedOut                       KwPrinterStateReasons = "timed-out"
	KwPrinterStateTonerEmpty                     KwPrinterStateReasons = "toner-empty"
	KwPrinterStateTonerLow                       KwPrinterStateReasons = "toner-low"

	// "printer-state-reasons" may also have one of the following
	// standard suffixes, indicating its level of severity:
	KwPrinterStateReport  KwPrinterStateReasons = "-report"
	KwPrinterStateWarning KwPrinterStateReasons = "-warning"
	KwPrinterStateError   KwPrinterStateReasons = "-error"
)

// Split splits KwPrinterStateReasons string into reason and severity suffix
// (which can be one of "-report", "-warning" or "-error".
//
// If string contains one of the known suffixes, reason itself and
// suffix will be returned separately. Otherwise, the first returned
// value will be unmodified string, and second will be "":
//
//	"media-low-warning"  ->  "media-low", "-warning"
//	"media-jam-error"    ->  "media-jam", "-error"
//	"shutdown"           ->  "shutdown",  ""
func (s KwPrinterStateReasons) Split() (
	reason, severity KwPrinterStateReasons) {

	idx := strings.LastIndexByte(string(s), '-')

	if idx >= 0 {
		prefix, suffix := s[:idx], s[idx:]
		switch suffix {
		case KwPrinterStateReport,
			KwPrinterStateWarning,
			KwPrinterStateError:
			return prefix, suffix
		}
	}

	return s, ""
}

// Reason splits KwPrinterStateReasons into reason and severity
// and returns its reason part.
func (s KwPrinterStateReasons) Reason() KwPrinterStateReasons {
	reason, _ := s.Split()
	return reason
}

// Severity splits KwPrinterStateReasons into reason and severity
// and returns its severity part.
func (s KwPrinterStateReasons) Severity() KwPrinterStateReasons {
	_, severity := s.Split()
	return severity
}

// KwSides represents standard keyword values for
// "sides" attribute.
//
// See RFC8011, 5.2.8.
type KwSides string

const (
	// KwSidesOneSided imposes one-side output
	KwSidesOneSided KwSides = "one-sided"

	// KwSidesTwoSidedLongEdge imposed two-sided output with
	// Impressions orientation suitable for the long edge binding
	KwSidesTwoSidedLongEdge KwSides = "two-sided-long-edge"

	// KwSidesTwoSidedShortEdge imposed two-sided output with
	// Impressions orientation suitable for the short edge binding
	KwSidesTwoSidedShortEdge KwSides = "two-sided-short-edge"
)

// KwURIAuthentication represents standard keyword values for
// "uri-authentication-supported" attribute.
//
// See RFC8011, 5.4.2.
type KwURIAuthentication string

const (
	// KwURIAuthenticationNone means that there is no
	// authentication mechanism associated with the URI.
	KwURIAuthenticationNone KwURIAuthentication = "none"

	// KwURIAuthenticationRequestingUserName means that Client
	// sends authenticated user by the "requesting-user-name"
	// operation attribute.
	KwURIAuthenticationRequestingUserName KwURIAuthentication = "requesting-user-name"

	// KwURIAuthenticationBasic means HTTP basic authentication.
	KwURIAuthenticationBasic KwURIAuthentication = "basic"

	// KwURIAuthenticationDigest means HTTP digest authentication.
	KwURIAuthenticationDigest KwURIAuthentication = "digest"

	// KwURIAuthenticationCertificate means TLS authentication
	// based on X.509 certificates.
	KwURIAuthenticationCertificate KwURIAuthentication = "certificate"
)

// KwURISecurity represents standard keyword values for
// "uri-security-supported" attribute.
//
// See RFC8011, 5.4.3.
type KwURISecurity string

const (
	// KwURISecurityNone means that there is no secure communication
	// channel in use for given URI
	KwURISecurityNone KwURISecurity = "none"

	// KwURISecurityTLS indicates TLS security
	KwURISecurityTLS KwURISecurity = "tls"
)

// KwWhichJobs represents standard keyword values for
// "which-jobs" attribute.
//
// RFC8011: 4.2.6.1.
// PWG5100.7: 8.5.
// PWG5100.11: 11.2.
// PWG5100.18: 9.8.
type KwWhichJobs string

const (
	// ----- RFC8011 -----

	// KwWhichJobsCompleted means all completed jobs, i.e. jobs
	// whose state is 'completed', 'canceled', or 'aborted'.
	KwWhichJobsCompleted KwWhichJobs = "completed"

	// KwWhichJobsNotCompleted means all non-completed jobs, i.e.
	// Job whose state is 'pending', 'processing', 'processing-stopped',
	// or 'pending-held'.
	KwWhichJobsNotCompleted KwWhichJobs = "not-completed"

	// ----- PWG5100.7 -----

	// KwWhichJobsAborted means all jobs in the 'aborted' state.
	KwWhichJobsAborted KwWhichJobs = "aborted"

	// KwWhichJobsAll means all jobs regardless of state.
	KwWhichJobsAll KwWhichJobs = "all"

	// KwWhichJobsCanceled means all jobs in the 'canceled' state.
	KwWhichJobsCanceled KwWhichJobs = "canceled"

	// KwWhichJobsPending means all jobs in the 'pending' state.
	KwWhichJobsPending KwWhichJobs = "pending"

	// KwWhichJobsPendingHeld means all jobs in the 'pending-held' state.
	KwWhichJobsPendingHeld KwWhichJobs = "pending-held"

	// KwWhichJobsProcessing means all jobs in the 'processing' state.
	KwWhichJobsProcessing KwWhichJobs = "processing"

	// KwWhichJobsProcessinStopped means all jobs in the
	// 'processing-stopped' state
	KwWhichJobsProcessinStopped KwWhichJobs = "processing-stopped"

	// ----- PWG5100.11 -----

	// KwWhichJobsProofPrint means all jobs that have been submitted
	// using the "proof-print" Job Template attribute and which are in
	// the ‘completed’, ‘canceled’, or ‘aborted’ state.
	KwWhichJobsProofPrint KwWhichJobs = "proof-print"

	// KwWhichJobsSaved means all jobs that have been saved using the
	// "job-save-disposition" Job Template attribute and which are in
	// the ‘completed’, ‘canceled’, or ‘aborted’ state.
	KwWhichJobsSaved KwWhichJobs = "saved"

	// ----- PWG5100.18 -----

	// KwWhichJobsFetchable means those jobs whose "job-state-reasons"
	// Job Description attribute contains the value 'job-fetchable' are to
	// be returned by the Get-Jobs operation.
	KwWhichJobsFetchable KwWhichJobs = "fetchable"
)

// KwRequestedAttribute represents a keyword value of the
// "requested-attributes" operation attribute.
type KwRequestedAttribute string

const (
	// KwRequestedAttributeAll requests all Job attribute groups.
	KwRequestedAttributeAll KwRequestedAttribute = "all"

	// KwRequestedAttributeJobDescription requests the job-description
	// attribute group.
	KwRequestedAttributeJobDescription KwRequestedAttribute = "job-description"

	// KwRequestedAttributeJobTemplate requests the job-template
	// attribute group.
	KwRequestedAttributeJobTemplate KwRequestedAttribute = "job-template"

	// KwRequestedAttributeJobID is the default "job-id" attribute
	// returned by Get-Jobs when requested-attributes is omitted.
	KwRequestedAttributeJobID KwRequestedAttribute = "job-id"

	// KwRequestedAttributeJobURI is the default "job-uri" attribute
	// returned by Get-Jobs when requested-attributes is omitted.
	KwRequestedAttributeJobURI KwRequestedAttribute = "job-uri"
)

// kwRegisteredTypes lists all registered keyword types for IPP codec.
var kwRegisteredTypes = map[reflect.Type]struct{}{
	// Types, defined here
	reflect.TypeOf(KwCompression("")):              struct{}{},
	reflect.TypeOf(KwJobDelayOutputUntil("")):      struct{}{},
	reflect.TypeOf(KwJobHoldUntil("")):             struct{}{},
	reflect.TypeOf(KwJobSheets("")):                struct{}{},
	reflect.TypeOf(KwJobSpooling("")):              struct{}{},
	reflect.TypeOf(KwJobStateReasons("")):          struct{}{},
	reflect.TypeOf(KwMediaBackCoating("")):         struct{}{},
	reflect.TypeOf(KwMultipleDocumentHandling("")): struct{}{},
	reflect.TypeOf(KwPdlOverride("")):              struct{}{},
	reflect.TypeOf(KwPrinterStateReasons("")):      struct{}{},
	reflect.TypeOf(KwSides("")):                    struct{}{},
	reflect.TypeOf(KwURIAuthentication("")):        struct{}{},
	reflect.TypeOf(KwURISecurity("")):              struct{}{},
	reflect.TypeOf(KwWhichJobs("")):                struct{}{},
	reflect.TypeOf(KwRequestedAttribute("")):       struct{}{},

	// Types, defined at separate source files
	reflect.TypeOf(KwColor("")):       struct{}{},
	reflect.TypeOf(KwDeviceClass("")): struct{}{},
	reflect.TypeOf(KwMedia("")):       struct{}{},

	// Input attribute keyword types, defined in kwinputattr.go
	reflect.TypeOf(KwInputSource("")):       struct{}{},
	reflect.TypeOf(KwInputColorMode("")):    struct{}{},
	reflect.TypeOf(KwInputContentType("")):  struct{}{},
	reflect.TypeOf(KwInputFilmScanMode("")): struct{}{},
}
