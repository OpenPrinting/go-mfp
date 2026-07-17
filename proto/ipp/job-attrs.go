// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Job and Job Template Attributes

package ipp

import (
	"time"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/proto/ipp/iana"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// JobCreateOperation contains operation attributes common for
// the job creation requests.
type JobCreateOperation struct {
	OperationGroup

	PrinterURI              string               `ipp:"printer-uri"`
	RequestingUserName      optional.Val[string] `ipp:"requesting-user-name"`
	Compression             optional.Val[string] `ipp:"compression"`
	DocumentFormat          optional.Val[string] `ipp:"document-format"`
	DocumentName            optional.Val[string] `ipp:"document-name"`
	DocumentNaturalLanguage optional.Val[string] `ipp:"document-natural-language"`
	IppAttributeFidelity    optional.Val[bool]   `ipp:"ipp-attribute-fidelity"`
	JobImpressions          optional.Val[int]    `ipp:"job-impressions"`
	JobKOctets              optional.Val[int]    `ipp:"job-k-octets"`
	JobMediaSheets          optional.Val[int]    `ipp:"job-media-sheets"`
	JobName                 optional.Val[string] `ipp:"job-name"`
	RequestingUserURI       optional.Val[string] `ipp:"requesting-user-uri"`

	// PWG5100.17, 7.1.1: scan-job operation attributes.
	CompressionAccepted    []KwCompression                `ipp:"compression-accepted"`
	DocumentFormatAccepted []string                       `ipp:"document-format-accepted"`
	InputAttributes        optional.Val[InputAttributes]  `ipp:"input-attributes"`
	OutputAttributes       optional.Val[OutputAttributes] `ipp:"output-attributes"`
}

// ToAbstract converts [JobCreateOperation] scan parameters into an
// [abstract.ScannerRequest].
//
// The returned request is partially populated — only fields that have
// an IPP representation are set. Callers should pass the result through
// [abstract.ScannerCapabilities.FillRequest] to validate and fill in
// any remaining defaults.
func (op *JobCreateOperation) ToAbstract() abstract.ScannerRequest {
	req := abstract.ScannerRequest{}

	if op.DocumentFormat != nil {
		req.DocumentFormat = optional.Get(op.DocumentFormat)
	} else if len(op.DocumentFormatAccepted) > 0 {
		req.DocumentFormat = op.DocumentFormatAccepted[0]
	}

	if op.InputAttributes == nil {
		return req
	}
	inp := *op.InputAttributes

	if inp.InputSource != nil {
		switch optional.Get(inp.InputSource) {
		case KwInputSourcePlaten:
			req.Input = abstract.InputPlaten
		case KwInputSourceADF:
			req.Input = abstract.InputADF
			req.ADFMode = abstract.ADFModeSimplex
			if inp.InputSides != nil &&
				optional.Get(inp.InputSides) == KwSidesTwoSidedLongEdge {
				req.ADFMode = abstract.ADFModeDuplex
			}
		}
	}

	if inp.InputColorMode != nil {
		req.ColorMode, req.ColorDepth = inputColorModeToAbstract(
			optional.Get(inp.InputColorMode))
	}

	if inp.InputResolution != nil {
		r := optional.Get(inp.InputResolution)
		req.Resolution = abstract.Resolution{
			XResolution: r.Xres,
			YResolution: r.Yres,
		}
	}

	if len(inp.InputScanRegions) > 0 {
		reg := inp.InputScanRegions[0]
		if reg.XOrigin != nil {
			req.Region.XOffset = abstract.Dimension(optional.Get(reg.XOrigin))
		}
		if reg.YOrigin != nil {
			req.Region.YOffset = abstract.Dimension(optional.Get(reg.YOrigin))
		}
		if reg.XDimension != nil {
			req.Region.Width = abstract.Dimension(optional.Get(reg.XDimension))
		}
		if reg.YDimension != nil {
			req.Region.Height = abstract.Dimension(optional.Get(reg.YDimension))
		}
	}

	// Brightness, Contrast, Sharpness
	req.Brightness = inp.InputBrightness
	req.Contrast = inp.InputContrast
	req.Sharpen = inp.InputSharpness

	// OutputAttributes → NoiseRemoval
	if op.OutputAttributes != nil {
		req.NoiseRemoval = optional.Get(op.OutputAttributes).NoiseRemoval
	}

	return req
}

// JobDescriptionAttrs contains IANA "job-description" attributes —
// static identity fields set once at creation that never change.
type JobDescriptionAttrs struct {
	JobDescriptionGroup

	JobID                  int                  `ipp:"job-id"`
	JobName                optional.Val[string] `ipp:"job-name"`
	JobOriginatingUserName optional.Val[string] `ipp:"job-originating-user-name"`
	JobURI                 string               `ipp:"job-uri"`
}

// JobStatusAttrs contains IANA "job-status" attributes —
// dynamic state fields updated throughout the job lifecycle.
type JobStatusAttrs struct {
	JobStatusGroup

	JobImpressionsCompleted optional.Val[int]    `ipp:"job-impressions-completed"`
	JobMediaSheetsCompleted optional.Val[int]    `ipp:"job-media-sheets-completed"`
	JobState                EnJobState           `ipp:"job-state"`
	JobStateMessage         optional.Val[string] `ipp:"job-state-message"`
	JobStateReasons         []KwJobStateReasons  `ipp:"job-state-reasons"`
	NumberOfInterveningJobs optional.Val[int]    `ipp:"number-of-intervening-jobs"`
}

type JobDescriptionAndStatus struct {
	ObjectRawAttrs
	JobDescriptionAttrs
	JobStatusAttrs
}

// JobGroupEntry holds the complete set of attributes returned for a single
// Job object in a Get-Jobs or Get-Job-Attributes response.
//
// On the wire, job-description, job-status and job-template attributes all
// arrive in a single flat Job group, so they are combined here into one
// flat structure decoded and encoded in a single pass.
type JobGroupEntry struct {
	ObjectRawAttrs
	JobDescriptionAttrs
	JobStatusAttrs
	JobTemplateAttrs
}

// DecodeJobDescriptionAndStatus decodes [JobDescriptionAndStatus] from
// [goipp.Attributes].
func DecodeJobDescriptionAndStatus(attrs goipp.Attributes, opt *DecoderOptions) (
	*JobDescriptionAndStatus, error) {

	job := &JobDescriptionAndStatus{}
	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(job, attrs)
	if err != nil {
		return nil, err
	}
	return job, nil
}

type JobTemplate struct {
	ObjectRawAttrs
	JobTemplateAttrs
}

// JobTemplateAttrs contains the IANA "job-template" attributes only.
type JobTemplateAttrs struct {
	JobTemplateGroup

	// RFC8011, Internet Printing Protocol/1.1: Model and Semantics
	// 5.2 Job Template Attributes
	Copies                   optional.Val[int]                        `ipp:"copies"`
	Finishings               []int                                    `ipp:"finishings"`
	JobHoldUntil             optional.Val[KwJobHoldUntil]             `ipp:"job-hold-until"`
	JobPriority              optional.Val[int]                        `ipp:"job-priority"`
	JobSheets                optional.Val[KwJobSheets]                `ipp:"job-sheets"`
	Media                    optional.Val[KwMedia]                    `ipp:"media"`
	MultipleDocumentHandling optional.Val[KwMultipleDocumentHandling] `ipp:"multiple-document-handling"`
	NumberUp                 optional.Val[int]                        `ipp:"number-up"`
	OrientationRequested     optional.Val[int]                        `ipp:"orientation-requested"`
	PageRanges               []goipp.Range                            `ipp:"page-ranges"`
	PrinterResolution        optional.Val[goipp.Resolution]           `ipp:"printer-resolution"`
	PrintQuality             optional.Val[int]                        `ipp:"print-quality"`
	Sides                    optional.Val[KwSides]                    `ipp:"sides"`

	// PWG5100.2: IPP “output-bin” attribute extension
	OutputBin optional.Val[string] `ipp:"output-bin"`

	// PWG5100.7: IPP Job Extensions v2.1 (JOBEXT)
	// 6.8 Job Template Attributes
	JobDelayOutputUntil     optional.Val[KwJobDelayOutputUntil] `ipp:"job-delay-output-until"`
	JobDelayOutputUntilTime optional.Val[time.Time]             `ipp:"job-delay-output-until-time"`
	JobHoldUntilTime        optional.Val[time.Time]             `ipp:"job-hold-until-time"`
	JobAccountID            optional.Val[string]                `ipp:"job-account-id"`
	JobAccountingUserID     optional.Val[string]                `ipp:"job-accounting-user-id"`
	JobCancelAfter          optional.Val[int]                   `ipp:"job-cancel-after"`
	JobRetainUntil          optional.Val[string]                `ipp:"job-retain-until"`
	JobRetainUntilInterval  optional.Val[int]                   `ipp:"job-retain-until-interval"`
	JobRetainUntilTime      optional.Val[time.Time]             `ipp:"job-retain-until-time"`
	JobSheetMessage         optional.Val[string]                `ipp:"job-sheet-message"`
	JobSheetsCol            JobSheets                           `ipp:"job-sheets-col"`
	PrintContentOptimize    optional.Val[string]                `ipp:"print-content-optimize"`

	// PWG5100.11: IPP Job and Printer Extensions – Set 2 (JPS2)
	// 7 Job Template Attributes
	FeedOrientation  optional.Val[string] `ipp:"feed-orientation"`
	JobPhoneNumber   optional.Val[string] `ipp:"job-phone-number"`
	JobRecipientName optional.Val[string] `ipp:"job-recipient-name"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.2 Job and Document Template Attributes
	JobErrorAction       optional.Val[string] `ipp:"job-error-action"`
	MediaOverprint       MediaOverprint       `ipp:"media-overprint"`
	PrintColorMode       optional.Val[string] `ipp:"print-color-mode"`
	PrintRenderingIntent optional.Val[string] `ipp:"print-rendering-intent"`
	PrintScaling         optional.Val[string] `ipp:"print-scaling"`

	// Wi-Fi Peer-to-Peer Services Print (P2Ps-Print)
	// Technical Specification
	// (for Wi-Fi Direct® services certification)
	PclmSourceResolution optional.Val[goipp.Resolution] `ipp:"pclm-source-resolution"`
}

// DecodeJobTemplate decodes [JobTemplate] from
// [goipp.Attributes].
func DecodeJobTemplate(attrs goipp.Attributes, opt *DecoderOptions) (
	*JobTemplate, error) {

	job := &JobTemplate{}
	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(job, attrs)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// DecodeJobGroupEntry decodes a single Job attribute group (as returned
// in Get-Jobs and Get-Job-Attributes responses) into a [JobGroupEntry].
//
// The whole flat group (job-description, job-status and job-template
// attributes) is decoded in a single pass.
func DecodeJobGroupEntry(attrs goipp.Attributes, opt *DecoderOptions) (
	JobGroupEntry, error) {

	var entry JobGroupEntry
	dec := NewDecoder(opt)
	defer dec.Free()

	if err := dec.Decode(&entry, attrs); err != nil {
		return JobGroupEntry{}, err
	}

	return entry, nil
}

// JobTemplateCapabilities are attributes, included into the Printer
// Description and describing possible settings for a [JobTemplate]
// (the "*-default", "*-supported" and "*-ready" values).
type JobTemplateCapabilities struct {
	// RFC8011, Internet Printing Protocol/1.1: Model and Semantics
	// 5.2 Job Template Attributes
	CopiesDefault                     optional.Val[int]                        `ipp:"copies-default"`
	CopiesSupported                   optional.Val[goipp.Range]                `ipp:"copies-supported"`
	FinishingsDefault                 []int                                    `ipp:"finishings-default"`
	FinishingsSupported               []int                                    `ipp:"finishings-supported"`
	JobHoldUntilDefault               optional.Val[KwJobHoldUntil]             `ipp:"job-hold-until-default"`
	JobHoldUntilSupported             []KwJobHoldUntil                         `ipp:"job-hold-until-supported"`
	JobPriorityDefault                optional.Val[int]                        `ipp:"job-priority-default"`
	JobPrioritySupported              optional.Val[int]                        `ipp:"job-priority-supported"`
	JobSheetsDefault                  []KwJobSheets                            `ipp:"job-sheets-default"`
	JobSheetsSupported                []KwJobSheets                            `ipp:"job-sheets-supported"`
	MediaDefault                      optional.Val[KwMedia]                    `ipp:"media-default"`
	MediaReady                        []KwMedia                                `ipp:"media-ready"`
	MediaSupported                    []KwMedia                                `ipp:"media-supported"`
	MultipleDocumentHandlingDefault   optional.Val[KwMultipleDocumentHandling] `ipp:"multiple-document-handling-default"`
	MultipleDocumentHandlingSupported []KwMultipleDocumentHandling             `ipp:"multiple-document-handling-supported"`
	NumberUpDefault                   optional.Val[int]                        `ipp:"number-up-default"`
	NumberUpSupported                 []goipp.IntegerOrRange                   `ipp:"number-up-supported"`
	OrientationRequestedDefault       optional.Val[int]                        `ipp:"orientation-requested-default"`
	OrientationRequestedSupported     []int                                    `ipp:"orientation-requested-supported"`
	PageRangesSupported               optional.Val[bool]                       `ipp:"page-ranges-supported"`
	PrinterResolutionDefault          optional.Val[goipp.Resolution]           `ipp:"printer-resolution-default"`
	PrinterResolutionSupported        []goipp.Resolution                       `ipp:"printer-resolution-supported"`
	PrintQualityDefault               optional.Val[int]                        `ipp:"print-quality-default"`
	PrintQualitySupported             []int                                    `ipp:"print-quality-supported"`
	SidesDefault                      optional.Val[KwSides]                    `ipp:"sides-default"`
	SidesSupported                    []KwSides                                `ipp:"sides-supported"`

	// PWG5100.2: IPP “output-bin” attribute extension
	OutputBinDefault   optional.Val[string] `ipp:"output-bin-default"`
	OutputBinSupported []string             `ipp:"output-bin-supported"`

	// PWG5100.7: IPP Job Extensions v2.1 (JOBEXT)
	// 6.9 Printer Description Attributes
	JobAccountIDDefault              optional.Val[string]                `ipp:"job-account-id-default"`
	JobAccountIDSupported            optional.Val[bool]                  `ipp:"job-account-id-supported"`
	JobAccountingUserIDDefault       optional.Val[string]                `ipp:"job-accounting-user-id-default"`
	JobAccountingUserIDSupported     optional.Val[bool]                  `ipp:"job-accounting-user-id-supported"`
	JobCancelAfterDefault            optional.Val[int]                   `ipp:"job-cancel-after-default"`
	JobCancelAfterSupported          optional.Val[goipp.Range]           `ipp:"job-cancel-after-supported"`
	JobDelayOutputUntilDefault       optional.Val[KwJobDelayOutputUntil] `ipp:"job-delay-output-until-default"`
	JobDelayOutputUntilSupported     []KwJobDelayOutputUntil             `ipp:"job-delay-output-until-supported"`
	JobDelayOutputUntilTimeSupported optional.Val[goipp.Range]           `ipp:"job-delay-output-until-time-supported"`
	JobHoldUntilTimeSupported        optional.Val[bool]                  `ipp:"job-hold-until-time-supported"`
	JobRetainUntilDefault            optional.Val[string]                `ipp:"job-retain-until-default"`
	JobRetainUntilIntervalDefault    optional.Val[int]                   `ipp:"job-retain-until-interval-default"`
	JobRetainUntilIntervalSupported  optional.Val[goipp.Range]           `ipp:"job-retain-until-interval-supported"`
	JobRetainUntilSupported          []string                            `ipp:"job-retain-until-supported"`
	JobRetainUntilTimeSupported      optional.Val[goipp.Range]           `ipp:"job-retain-until-time-supported"`
	JobSheetsColDefault              optional.Val[JobSheets]             `ipp:"job-sheets-col-default"`
	JobSheetsColSupported            []string                            `ipp:"job-sheets-col-supported"`
	PrintContentOptimizeDefault      optional.Val[string]                `ipp:"print-content-optimize-default"`
	PrintContentOptimizeSupported    []string                            `ipp:"print-content-optimize-supported"`

	// PWG5100.11: IPP Job and Printer Extensions – Set 2 (JPS2)
	// 7 Job Template Attributes
	FeedOrientationDefault               optional.Val[string] `ipp:"feed-orientation-default"`
	FeedOrientationSupported             string               `ipp:"feed-orientation-supported"`
	JobPhoneNumberDefault                optional.Val[string] `ipp:"job-phone-number-default"`
	JobPhoneNumberSupported              optional.Val[bool]   `ipp:"job-phone-number-supported"`
	JobRecipientNameDefault              optional.Val[string] `ipp:"job-recipient-name-default"`
	JobRecipientNameSupported            optional.Val[bool]   `ipp:"job-recipient-name-supported"`
	PdlInitFileEntrySupported            []string             `ipp:"pdl-init-file-entry-supported"`
	PdlInitFileNameSubdirectorySupported optional.Val[bool]   `ipp:"pdl-init-file-name-subdirectory-supported"`
	PdlInitFileNameSupported             []string             `ipp:"pdl-init-file-name-supported"`
	PdlInitFileSupported                 []string             `ipp:"pdl-init-file-supported"`
	PrintProcessingAttributesSupported   []string             `ipp:"print-processing-attributes-supported"`
	SaveDispositionSupported             []string             `ipp:"save-disposition-supported"`
	SaveDocumentFormatDefault            optional.Val[string] `ipp:"save-document-format-default"`
	SaveDocumentFormatSupported          []string             `ipp:"save-document-format-supported"`
	SaveLocationDefault                  optional.Val[string] `ipp:"save-location-default"`
	SaveLocationSupported                []string             `ipp:"save-location-supported"`
	SaveNameSubdirectorySupported        optional.Val[bool]   `ipp:"save-name-subdirectory-supported"`
	SaveNameSupported                    optional.Val[bool]   `ipp:"save-name-supported"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.2 Job and Document Template Attributes
	// 6.5 Printer Description Attributes
	JobErrorActionDefault           optional.Val[string]         `ipp:"job-error-action-default"`
	JobErrorActionSupported         []string                     `ipp:"job-error-action-supported"`
	MediaOverprintDefault           optional.Val[MediaOverprint] `ipp:"media-overprint-default"`
	MediaOverprintDistanceSupported optional.Val[goipp.Range]    `ipp:"media-overprint-distance-supported"`
	MediaOverprintMethodSupported   []string                     `ipp:"media-overprint-method-supported"`
	MediaOverprintSupported         []string                     `ipp:"media-overprint-supported"`
	PrintColorModeDefault           optional.Val[string]         `ipp:"print-color-mode-default"`
	PrintColorModeSupported         []string                     `ipp:"print-color-mode-supported"`
	PrinterMandatoryJobAttributes   []string                     `ipp:"printer-mandatory-job-attributes"`
	PrintRenderingIntentDefault     optional.Val[string]         `ipp:"print-rendering-intent-default"`
	PrintRenderingIntentSupported   []string                     `ipp:"print-rendering-intent-supported"`
	PrintScalingDefault             optional.Val[string]         `ipp:"print-scaling-default"`
	PrintScalingSupported           []string                     `ipp:"print-scaling-supported"`

	// Wi-Fi Peer-to-Peer Services Print (P2Ps-Print)
	// Technical Specification
	// (for Wi-Fi Direct® services certification)
	PclmSourceResolutionDefault   optional.Val[goipp.Resolution] `ipp:"pclm-source-resolution-default"`
	PclmSourceResolutionSupported []goipp.Resolution             `ipp:"pclm-source-resolution-supported"`
}

// JobSheets represents "job-sheets-col" collection entry in
// JobTemplate
type JobSheets struct {
	JobSheets KwJobSheets `ipp:"job-sheets"`
	Media     string      `ipp:"media"`
	MediaCol  MediaCol    `ipp:"media-col"`
}

// JobPresets represents "job-presets-supported" collection entry
// in PrinterDescription
type JobPresets struct {
	PresetCategory string `ipp:"preset-category"`
	PresetName     string `ipp:"preset-name"`
	JobTemplate
}

// jobAttrGroups maps the standard attribute-group keywords used by
// Get-Jobs and Get-Job-Attributes to individual attribute names.
//
// See RFC8011, 4.3.4.
var jobAttrGroups = buildJobAttrGroups()

// jobAttributesFilter filters job attributes against the requested
// groups/names for Get-Jobs and Get-Job-Attributes responses.
var jobAttributesFilter = &filterAttributes{
	groups:   jobAttrGroups,
	standard: jobAttrGroups["all"],
}

func buildJobAttrGroups() map[string]generic.Set[string] {
	jobDescription := generic.NewSet[string]()
	for name := range iana.JobDescription {
		jobDescription.Add(name)
	}
	for name := range iana.JobStatus {
		jobDescription.Add(name)
	}

	template := generic.NewSet[string]()
	for name := range iana.JobTemplate {
		template.Add(name)
	}

	all := jobDescription.Clone()
	all.Merge(template)

	return map[string]generic.Set[string]{
		"all":             all,
		"job-description": jobDescription,
		"job-template":    template,
	}
}
