// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Job and Job Template Attributes

package ipp

import (
	"time"

	"github.com/OpenPrinting/goipp"
)

// JobAttributes are attributes, supplied with Job creation request
type JobAttributes struct {
	// RFC8011, Internet Printing Protocol/1.1: Model and Semantics
	// 5.2 Job Template Attributes
	Copies                   int                        `ipp:"?copies,>0"`
	Finishings               []int                      `ipp:"?finishings,enum"`
	JobHoldUntil             KwJobHoldUntil             `ipp:"?job-hold-until"`
	JobPriority              int                        `ipp:"?job-priority,1:100"`
	JobSheets                KwJobSheets                `ipp:"?job-sheets"`
	Media                    KwMedia                    `ipp:"?media"`
	MultipleDocumentHandling KwMultipleDocumentHandling `ipp:"?multiple-document-handling"`
	NumberUp                 int                        `ipp:"?number-up,>0"`
	OrientationRequested     int                        `ipp:"?orientation-requested,enum"`
	PageRanges               []goipp.IntegerOrRange     `ipp:"?page-ranges"`
	PrinterResolution        goipp.Resolution           `ipp:"?printer-resolution"`
	PrintQuality             int                        `ipp:"?print-quality,enum"`
	Sides                    KwSides                    `ipp:"?sides"`

	// PWG5100.7: IPP Job Extensions v2.1 (JOBEXT)
	// 6.8 Job Template Attributes
	JobDelayOutputUntil     KwJobDelayOutputUntil `ipp:"?job-delay-output-until"`
	JobDelayOutputUntilTime time.Time             `ipp:"?job-delay-output-until-time"`
	JobHoldUntilTime        time.Time             `ipp:"?job-hold-until-time"`
	JobAccountID            string                `ipp:"?job-account-id,name"`
	JobAccountingUserID     string                `ipp:"?job-accounting-user-id,name"`
	JobCancelAfter          int                   `ipp:"?job-cancel-after,0:MAX"`
	JobRetainUntil          string                `ipp:"?job-retain-until,keyword"`
	JobRetainUntilInterval  int                   `ipp:"?job-retain-until-interval,0:MAX"`
	JobRetainUntilTime      time.Time             `ipp:"?job-retain-until-time"`
	JobSheetMessage         string                `ipp:"?job-sheet-message,text"`
	JobSheetsCol            []JobSheets           `ipp:"?job-sheets-col"`

	// PWG5100.11: IPP Job and Printer Extensions – Set 2 (JPS2)
	// 7 Job Template Attributes
	FeedOrientation    string               `ipp:"?feed-orientation,keyword"`
	FontNameRequested  string               `ipp:"?font-name-requested,name"`
	FontSizeRequested  int                  `ipp:"?font-size-requested,>0"`
	JobPhoneNumber     string               `ipp:"?job-phone-number,uri"`
	JobRecipientName   string               `ipp:"?job-recipient-name,name"`
	JobSaveDisposition []JobSaveDisposition `ipp:"?job-save-disposition"`
	PdlInitFile        []JobPdlInitFile     `ipp:"?pdl-init-file"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.2 Job and Document Template Attributes
	JobErrorAction       string              `ipp:"?job-error-action,keyword"`
	MediaOverprint       []JobMediaOverprint `ipp:"?media-overprint"`
	PrintColorMode       string              `ipp:"?print-color-mode,keyword"`
	PrintRenderingIntent string              `ipp:"?print-rendering-intent,keyword"`
	PrintScaling         string              `ipp:"?print-scaling,keyword"`
}

// JobTemplate are attributes, included into the Printer Description and
// describing possible settings for JobAttributes
type JobTemplate struct {
	// RFC8011, Internet Printing Protocol/1.1: Model and Semantics
	// 5.2 Job Template Attributes
	CopiesDefault                     int                          `ipp:"?copies-default,>0"`
	CopiesSupported                   goipp.Range                  `ipp:"?copies-supported,>0"`
	FinishingsDefault                 []int                        `ipp:"?finishings-default,enum"`
	FinishingsSupported               []int                        `ipp:"?finishings-supported,enum"`
	JobHoldUntilDefault               KwJobHoldUntil               `ipp:"?job-hold-until-default"`
	JobHoldUntilSupported             []KwJobHoldUntil             `ipp:"?job-hold-until-supported"`
	JobPriorityDefault                int                          `ipp:"?job-priority-default,1:100"`
	JobPrioritySupported              int                          `ipp:"?job-priority-supported,1:100"`
	JobSheetsDefault                  KwJobSheets                  `ipp:"?job-sheets-default"`
	JobSheetsSupported                []KwJobSheets                `ipp:"?job-sheets-supported"`
	MediaDefault                      KwMedia                      `ipp:"?media-default"`
	MediaReady                        []KwMedia                    `ipp:"?media-ready"`
	MediaSupported                    []KwMedia                    `ipp:"?media-supported"`
	MultipleDocumentHandlingDefault   KwMultipleDocumentHandling   `ipp:"?multiple-document-handling-default"`
	MultipleDocumentHandlingSupported []KwMultipleDocumentHandling `ipp:"?multiple-document-handling-supported"`
	NumberUpDefault                   int                          `ipp:"?number-up-default,>0"`
	NumberUpSupported                 []goipp.IntegerOrRange       `ipp:"?number-up-supported,>0"`
	OrientationRequestedDefault       int                          `ipp:"?orientation-requested-default,enum"`
	OrientationRequestedSupported     []int                        `ipp:"?orientation-requested-supported,enum"`
	PageRangesSupported               bool                         `ipp:"?page-ranges-supported"`
	PrinterResolutionDefault          goipp.Resolution             `ipp:"?printer-resolution-default"`
	PrinterResolutionSupported        []goipp.Resolution           `ipp:"?printer-resolution-supported"`
	PrintQualityDefault               int                          `ipp:"?print-quality-default,enum"`
	PrintQualitySupported             []int                        `ipp:"?print-quality-supported,enum"`
	SidesDefault                      KwSides                      `ipp:"?sides-default"`
	SidesSupported                    []KwSides                    `ipp:"?sides-supported"`

	// PWG5100.7: IPP Job Extensions v2.1 (JOBEXT)
	// 6.9 Printer Description Attributes
	JobAccountIDDefault              string                  `ipp:"job-account-id-default,name|no-value"`
	JobAccountIDSupported            bool                    `ipp:"?job-account-id-supported"`
	JobAccountingUserIDDefault       string                  `ipp:"?job-accounting-user-id-default,name|no-value"`
	JobAccountingUserIDSupported     bool                    `ipp:"?job-accounting-user-id-supported"`
	JobCancelAfterDefault            int                     `ipp:"?job-cancel-after-default,0:MAX"`
	JobCancelAfterSupported          goipp.Range             `ipp:"?job-cancel-after-supported,0:MAX"`
	JobDelayOutputUntilDefault       KwJobDelayOutputUntil   `ipp:"?job-delay-output-until-default"`
	JobDelayOutputUntilSupported     []KwJobDelayOutputUntil `ipp:"?job-delay-output-until-supported"`
	JobDelayOutputUntilTimeSupported goipp.Range             `ipp:"?job-delay-output-until-time-supported,0:MAX"`
	JobHoldUntilTimeSupported        goipp.Range             `ipp:"?job-hold-until-time-supported,0:MAX"`
	JobRetainUntilDefault            string                  `ipp:"?job-retain-until-default,keyword"`
	JobRetainUntilIntervalDefault    int                     `ipp:"?job-retain-until-interval-default,0:MAX"`
	JobRetainUntilIntervalSupported  goipp.Range             `ipp:"?job-retain-until-interval-supported,0:MAX"`
	JobRetainUntilSupported          []string                `ipp:"?job-retain-until-supported,keyword"`
	JobRetainUntilTimeSupported      goipp.Range             `ipp:"?job-retain-until-time-supported,0:MAX"`
	JobSheetsColDefault              []JobSheets             `ipp:"?job-sheets-col-default,collection|no-value"`
	JobSheetsColSupported            []string                `ipp:"?job-sheets-col-supported,keyword"`

	// PWG5100.11: IPP Job and Printer Extensions – Set 2 (JPS2)
	// 7 Job Template Attributes
	FeedOrientationDefault               string               `ipp:"?feed-orientation-default,keyword"`
	FeedOrientationSupported             []string             `ipp:"?feed-orientation-supported,keyword"`
	FontNameRequestedDefault             string               `ipp:"?font-name-requested-default,name"`
	FontNameRequestedSupported           []string             `ipp:"?font-name-requested-supported,name"`
	FontSizeRequestedDefault             int                  `ipp:"?font-size-requested-default,>0"`
	FontSizeRequestedSupported           []int                `ipp:"?font-size-requested-supported,>0"`
	JobPhoneNumberDefault                string               `ipp:"?job-phone-number-default,uri"`
	JobPhoneNumberSupported              bool                 `ipp:"?job-phone-number-supported"`
	JobRecipientNameDefault              string               `ipp:"?job-recipient-name-default,name"`
	JobRecipientNameSupported            bool                 `ipp:"?job-recipient-name-supported"`
	JobSaveDispositionDefault            []JobSaveDisposition `ipp:"?job-save-disposition-default"`
	JobSaveDispositionSupported          []string             `ipp:"?job-save-disposition-supported,keyword"`
	PdlInitFileDefault                   []JobPdlInitFile     `ipp:"?pdl-init-file-default"`
	PdlInitFileEntrySupported            []string             `ipp:"?pdl-init-file-entry-supported,name"`
	PdlInitFileNameSubdirectorySupported bool                 `ipp:"?pdl-init-file-name-subdirectory-supported"`
	PdlInitFileNameSupported             []string             `ipp:"?pdl-init-file-name-supported,name"`
	PdlInitFileSupported                 []string             `ipp:"? pdl-init-file-supported,name"`
	PrintProcessingAttributesSupported   []string             `ipp:"?print-processing-attributes-supported,keyword"`
	SaveDispositionSupported             []string             `ipp:"?save-disposition-supported,keyword"`
	SaveDocumentFormatDefault            string               `ipp:"?save-document-format-default,mimeMediaType"`
	SaveDocumentFormatSupported          []string             `ipp:"?save-document-format-supported,mimeMediaType"`
	SaveInfoSupported                    []string             `ipp:"?save-info-supported,keyword"`
	SaveLocationDefault                  string               `ipp:"?save-location-default,uri"`
	SaveLocationSupported                []string             `ipp:"?save-location-supported,uri"`
	SaveNameSubdirectorySupported        bool                 `ipp:"?save-name-subdirectory-supported"`
	SaveNameSupported                    bool                 `ipp:"?save-name-supported"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.2 Job and Document Template Attributes
	// 6.5 Printer Description Attributes
	JobErrorActionDefault           string              `ipp:"?job-error-action-default,keyword"`
	JobErrorActionSupported         []string            `ipp:"?job-error-action-supported,keyword"`
	MediaOverprintDefault           []JobMediaOverprint `ipp:"?media-overprint-default"`
	MediaOverprintDistanceSupported goipp.Range         `ipp:"?media-overprint-distance-supported,0:MAX"`
	MediaOverprintMethodSupported   []string            `ipp:"?media-overprint-method-supported,keyword"`
	MediaOverprintSupported         []string            `ipp:"?media-overprint-supported,keyword"`
	PrintColorModeDefault           string              `ipp:"?print-color-mode-default,keyword"`
	PrintColorModeSupported         []string            `ipp:"?print-color-mode-supported,keyword"`
	PrinterMandatoryJobAttributes   []string            `ipp:"?printer-mandatory-job-attributes,keyword"`
	PrintRenderingIntentDefault     string              `ipp:"?print-rendering-intent-default,keyword"`
	PrintRenderingIntentSupported   []string            `ipp:"?print-rendering-intent-supported,keyword"`
	PrintScalingDefault             string              `ipp:"?print-scaling-default,keyword"`
	PrintScalingSupported           []string            `ipp:"?print-scaling-supported,keyword"`
}

// MediaCol is the "media-col", "media-col-xxx" collection entry.
// It is used in many places.
//
// PWG5100.3: 3.13., Table 10.
// PWG5100.7: 6.3.1., Table 6.
type MediaCol struct {
	// ----- PWG5100.3 -----
	KwMediaBackCoating KwMediaBackCoating `ipp:"?media-back-coating"`
	MediaColor         KwColor            `ipp:"?media-color"`
	MediaFrontCoating  KwMediaBackCoating `ipp:"?media-front-coating"`
	MediaHoleCount     int                `ipp:"?media-hole-count,0:MAX"`
	MediaInfo          string             `ipp:"?media-info,text"`
	MediaKey           KwMedia            `ipp:"?media-key"`
	MediaOrderCount    int                `ipp:"?media-order-count,1:MAX"`
	MediaPrePrinted    string             `ipp:"?media-pre-printed,keyword"`
	MediaRecycled      string             `ipp:"?media-recycled,keyword"`
	MediaSize          MediaSize          `ipp:"?media-size"`
	MediaType          string             `ipp:"?media-type,keyword"`
	MediaWeightMetric  int                `ipp:"?media-weight-metric,0:MAX"`

	// ----- PWG5100.7 -----
	MediaBottomMargin     int                   `ipp:"?media-bottom-margin,0:MAX"`
	MediaGrain            string                `ipp:"?media-grain,keyword"`
	MediaLeftMargin       int                   `ipp:"?media-left-margin,0:MAX"`
	MediaRightMargin      int                   `ipp:"?media-right-margin,0:MAX"`
	MediaSizeName         string                `ipp:"?media-size-name,keyword"`
	MediaSourceProperties MediaSourceProperties `ipp:"?media-source-properties"`
	MediaSource           string                `ipp:"?media-source,keyword"`
	MediaThickness        int                   `ipp:"?media-thickness,1:MAX"`
	MediaTooth            string                `ipp:"?media-tooth,keyword"`
	MediaTopMargin        int                   `ipp:"?media-top-margin,0:MAX"`
}

// MediaSize represents media size parameters (which may be either
// pair of integers or pair of ranges) and used in many places
type MediaSize struct {
	XDimension goipp.IntegerOrRange `ipp:"x-dimension,0:MAX"`
	YDimension goipp.IntegerOrRange `ipp:"y-dimension,0:MAX"`
}

// MediaSourceProperties represents "media-source-properties"
// collectiobn in MediaCol
type MediaSourceProperties struct {
	MediaSourceFeedDirection   string `ipp:"?media-source-feed-direction,keyword"`
	MediaSourceFeedOrientation int    `ipp:"?media-source-feed-orientation,enum"`
}

// JobSheets represents "job-sheets-col" collection entry in
// JobAttributes
type JobSheets struct {
	JobSheets KwJobSheets `ipp:"job-sheets"`
	Media     string      `ipp:"media,keyword"`
	MediaCol  []MediaCol  `ipp:"media-col"`
}

// JobSaveDisposition represents "job-save-disposition"
// collection entry in JobAttributes and "job-save-disposition-default"
// entry in JobTemplate
type JobSaveDisposition struct {
	SaveDisposition string        `ipp:"save-disposition,keyword"`
	SaveInfo        []JobSaveInfo `ipp:"?save-info"`
}

// JobSaveInfo represents "save-info" collection entry
// in JobSaveDisposition
type JobSaveInfo struct {
	SaveLocation       string `ipp:"?save-location,uri"`
	SaveName           string `ipp:"?save-name,name"`
	SaveDocumentFormat string `ipp:"?save-document-format,mimeMediaType"`
}

// JobPdlInitFile represents "pdl-init-file" collection entry
// in JobAttributes
type JobPdlInitFile struct {
	PdlInitFileLocation string `ipp:"?pdl-init-file-location,uri"`
	PdlInitFileName     string `ipp:"?pdl-init-file-name,name"`
	PdlInitFileEntry    string `ipp:"?pdl-init-file-entry,name"`
}

// JobMediaOverprint represents "media-overprint" collection entry
// in JobAttributes
type JobMediaOverprint struct {
	MediaOverprintDistance int    `ipp:"media-overprint-distance,0:MAX"`
	MediaOverprintMethod   string `ipp:"media-overprint-method,keyword"`
}

// JobPresets represents "job-presets-supported" collection entry
// in PrinterDescription
type JobPresets struct {
	PresetCategory string `ipp:"?preset-category,keyword"`
	PresetName     string `ipp:"?preset-name,name"`
	JobAttributes
}
