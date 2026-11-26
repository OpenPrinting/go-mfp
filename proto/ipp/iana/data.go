// MFP - Miulti-Function Printers and scanners toolkit
// IANA registrations for IPP
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP registrations database
//
// THIS IS GENERATED FILE. DON'T EDIT!

package iana

import (
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/goipp"
)

// CUPSDeviceAttributes is the CUPS Device Attributes attributes
var CUPSDeviceAttributes = map[string]*DefAttr{
	// CUPS Device Attributes/device-class (CUPS)
	"device-class": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// CUPS Device Attributes/device-id (CUPS)
	"device-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS Device Attributes/device-info (CUPS)
	"device-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS Device Attributes/device-location (CUPS)
	"device-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS Device Attributes/device-make-and-model (CUPS)
	"device-make-and-model": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS Device Attributes/device-uri (CUPS)
	"device-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
}

// CUPSPPDAttributes is the CUPS PPD Attributes attributes
var CUPSPPDAttributes = map[string]*DefAttr{
	// CUPS PPD Attributes/ppd-device-id (CUPS)
	"ppd-device-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS PPD Attributes/ppd-make (CUPS)
	"ppd-make": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS PPD Attributes/ppd-make-and-model (CUPS)
	"ppd-make-and-model": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS PPD Attributes/ppd-model-number (CUPS)
	"ppd-model-number": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// CUPS PPD Attributes/ppd-name (CUPS)
	"ppd-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// CUPS PPD Attributes/ppd-natural-language (CUPS)
	"ppd-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// CUPS PPD Attributes/ppd-product (CUPS)
	"ppd-product": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS PPD Attributes/ppd-psversion (CUPS)
	"ppd-psversion": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// CUPS PPD Attributes/ppd-type (CUPS)
	"ppd-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
}

// CUPSPrinterClassAttributes is the CUPS Printer Class Attributes attributes
var CUPSPrinterClassAttributes = map[string]*DefAttr{
	// CUPS Printer Class Attributes/member-names (CUPS)
	"member-names": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// CUPS Printer Class Attributes/member-uris (CUPS)
	"member-uris": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
}

// DocumentDescription is the Document Description attributes
var DocumentDescription = map[string]*DefAttr{
	// Document Description/document-name (PWG5100.5)
	"document-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
}

// DocumentStatus is the Document Status attributes
var DocumentStatus = map[string]*DefAttr{
	// Document Status/attributes-charset (PWG5100.5)
	"attributes-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Document Status/attributes-natural-language (PWG5100.5)
	"attributes-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Document Status/chamber-humidity-actual (PWG5100.21)
	"chamber-humidity-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/chamber-temperature-actual (PWG5100.21)
	"chamber-temperature-actual": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/compression (PWG5100.5)
	"compression": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/copies-actual (PWG5100.5)
	"copies-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/cover-back-actual (PWG5100.5)
	"cover-back-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/cover-front-actual (PWG5100.5)
	"cover-front-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/date-time-at-completed (PWG5100.5)
	"date-time-at-completed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Document Status/date-time-at-created (PWG5100.5)
	"date-time-at-created": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Document Status/date-time-at-creation (PWG5100.5)
	"date-time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Document Status/date-time-at-processing (PWG5100.5)
	"date-time-at-processing": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Document Status/detailed-status-messages (PWG5100.5)
	"detailed-status-messages": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-access-errors (PWG5100.5)
	"document-access-errors": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-charset (PWG5100.5)
	"document-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Document Status/document-digital-signature (PWG5100.5)
	"document-digital-signature": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/document-format (PWG5100.5)
	"document-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Document Status/document-format-details (PWG5100.7)
	"document-format-details": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/document-format-details-detected (PWG5100.7)
	"document-format-details-detected": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/document-format-detected (PWG5100.5)
	"document-format-detected": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Document Status/document-format-ready (PWG5100.18)
	"document-format-ready": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Document Status/document-format-supplied (PWG5100.7)
	"document-format-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Document Status/document-format-version (PWG5100.5)
	"document-format-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-format-version-detected (PWG5100.5)
	"document-format-version-detected": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-format-version-supplied (PWG5100.7)
	"document-format-version-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-job-id (PWG5100.5)
	"document-job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/document-job-uri (PWG5100.5)
	"document-job-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Document Status/document-message (PWG5100.5)
	"document-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-message-supplied (PWG5100.7)
	"document-message-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-metadata (PWG5100.13)
	"document-metadata": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Document Status/document-name-supplied (PWG5100.7)
	"document-name-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Document Status/document-natural-language (PWG5100.5)
	"document-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Document Status/document-number (PWG5100.5)
	"document-number": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/document-printer-uri (PWG5100.5)
	"document-printer-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Document Status/document-resource-ids (PWG5100.22)
	"document-resource-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/document-state (PWG5100.5)
	"document-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Status/document-state-message (PWG5100.5)
	"document-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/document-state-reasons (PWG5100.5)
	"document-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/document-uri (PWG5100.5)
	"document-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Document Status/document-uuid (PWG5100.13)
	"document-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Document Status/errors-count (PWG5100.7)
	"errors-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/finishings-actual (PWG5100.5)
	"finishings-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Status/finishings-col-actual (PWG5100.5)
	"finishings-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/force-front-side-actual (PWG5100.5)
	"force-front-side-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/imposition-template-actual (PWG5100.5)
	"imposition-template-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Status/impressions (PWG5100.5)
	"impressions": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/impressions-col (PWG5100.7)
	"impressions-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Status/impressions-col/blank (PWG5100.7)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/blank-two-sided (PWG5100.7)
			"blank-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/full-color-two-sided (PWG5100.7)
			"full-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/highlight-color (PWG5100.7)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/highlight-color-two-sided (PWG5100.7)
			"highlight-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/monochrome (PWG5100.7)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/impressions-col/monochrome-two-sided (PWG5100.7)
			"monochrome-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Document Status/impressions-completed (PWG5100.5)
	"impressions-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/impressions-completed-col (XEROX20150505)
	"impressions-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/input-attributes-actual (PWG5100.15)
	"input-attributes-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/insert-sheet-actual (PWG5100.5)
	"insert-sheet-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/k-octets (PWG5100.5)
	"k-octets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/k-octets-processed (PWG5100.5)
	"k-octets-processed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/last-document (PWG5100.5)
	"last-document": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Document Status/materials-col-actual (PWG5100.21)
	"materials-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/media-actual (PWG5100.5)
	"media-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Status/media-col-actual (PWG5100.5)
	"media-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/media-sheets (PWG5100.5)
	"media-sheets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/media-sheets-col (PWG5100.7)
	"media-sheets-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Status/media-sheets-col/blank (PWG5100.7)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/media-sheets-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/media-sheets-col/highlight-color (PWG5100.7)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/media-sheets-col/monochrome (XEROX20150505)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Document Status/media-sheets-completed (PWG5100.5)
	"media-sheets-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/media-sheets-completed-col (PWG5100.5)
	"media-sheets-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/more-info (PWG5100.5)
	"more-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Document Status/multiple-object-handling-actual (PWG5100.21)
	"multiple-object-handling-actual": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/number-up-actual (PWG5100.5)
	"number-up-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/orientation-requested-actual (PWG5100.5)
	"orientation-requested-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Status/output-attributes-actual (PWG5100.17)
	"output-attributes-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/output-bin-actual (PWG5100.5)
	"output-bin-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Document Status/output-device-actual (PWG5100.7)
	"output-device-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Document Status/output-device-assigned (PWG5100.5)
	"output-device-assigned": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Document Status/output-device-document-state (PWG5100.18)
	"output-device-document-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Status/output-device-document-state-message (PWG5100.18)
	"output-device-document-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Document Status/output-device-document-state-reasons (PWG5100.18)
	"output-device-document-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/page-delivery-actual (PWG5100.5)
	"page-delivery-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/page-order-received-actual (PWG5100.5)
	"page-order-received-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/page-ranges-actual (PWG5100.5)
	"page-ranges-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Document Status/pages (PWG5100.13)
	"pages": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/pages-col (PWG5100.7)
	"pages-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Status/pages-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Status/pages-col/monochrome (PWG5100.7)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Document Status/pages-completed (PWG5100.13)
	"pages-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/pages-completed-col (PWG5100.7)
	"pages-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/platform-temperature-actual (PWG5100.21)
	"platform-temperature-actual": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/presentation-direction-number-up-actual (PWG5100.5)
	"presentation-direction-number-up-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/print-accuracy-actual (PWG5100.21)
	"print-accuracy-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/print-base-actual (PWG5100.21)
	"print-base-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/print-color-mode-actual (PWG5100.13)
	"print-color-mode-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/print-content-optimize-actual (PWG5100.7)
	"print-content-optimize-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/print-objects-actual (PWG5100.21)
	"print-objects-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/print-quality-actual (PWG5100.5)
	"print-quality-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Status/print-supports-actual (PWG5100.21)
	"print-supports-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/printer-resolution-actual (PWG5100.5)
	"printer-resolution-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Document Status/printer-up-time (PWG5100.5)
	"printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/separator-sheets-actual (PWG5100.5)
	"separator-sheets-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Status/sheet-completed-copy-number (PWG5100.5)
	"sheet-completed-copy-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/sides-actual (PWG5100.5)
	"sides-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/time-at-completed (PWG5100.5)
	"time-at-completed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/time-at-creation (PWG5100.5)
	"time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/time-at-processing (PWG5100.5)
	"time-at-processing": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/warnings-count (PWG5100.7)
	"warnings-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/x-image-position-actual (PWG5100.5)
	"x-image-position-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/x-image-shift-actual (PWG5100.5)
	"x-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/x-side1-image-shift-actual (PWG5100.5)
	"x-side1-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/x-side2-image-shift-actual (PWG5100.5)
	"x-side2-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/y-image-position-actual (PWG5100.5)
	"y-image-position-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Status/y-image-shift-actual (PWG5100.5)
	"y-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/y-side1-image-shift-actual (PWG5100.5)
	"y-side1-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Status/y-side2-image-shift-actual (PWG5100.5)
	"y-side2-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
}

// DocumentTemplate is the Document Template attributes
var DocumentTemplate = map[string]*DefAttr{
	// Document Template/chamber-humidity (PWG5100.21)
	"chamber-humidity": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/chamber-temperature (PWG5100.21)
	"chamber-temperature": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/copies (PWG5100.5)
	"copies": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/cover-back (PWG5100.5)
	"cover-back": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Template/cover-front (PWG5100.5)
	"cover-front": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Template/feed-orientation (PWG5100.5)
	"feed-orientation": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/finishings (PWG5100.5)
	"finishings": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Template/finishings-col (PWG5100.5)
	"finishings-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Template/force-front-side (PWG5100.5)
	"force-front-side": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/imposition-template (PWG5100.5)
	"imposition-template": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Template/insert-sheet (PWG5100.5)
	"insert-sheet": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Template/materials-col (PWG5100.21)
	"materials-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Template/materials-col/material-amount (PWG5100.21)
			"material-amount": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-amount-units (PWG5100.21)
			"material-amount-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/materials-col/material-color (PWG5100.21)
			"material-color": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/materials-col/material-diameter (PWG5100.21)
			"material-diameter": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-diameter-tolerance (PWG5100.21)
			"material-diameter-tolerance": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-fill-density (PWG5100.21)
			"material-fill-density": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-key (PWG5100.21)
			"material-key": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/materials-col/material-name (PWG5100.21)
			"material-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Document Template/materials-col/material-nozzle-diameter (PWG5100.21)
			"material-nozzle-diameter": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-purpose (PWG5100.21)
			"material-purpose": &DefAttr{
				SetOf: true,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/materials-col/material-rate (PWG5100.21)
			"material-rate": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-rate-units (PWG5100.21)
			"material-rate-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/materials-col/material-retraction (PWG5100.21)
			"material-retraction": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Document Template/materials-col/material-shell-thickness (PWG5100.21)
			"material-shell-thickness": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/materials-col/material-temperature (PWG5100.21)
			"material-temperature": &DefAttr{
				SetOf: false,
				Min:   -273,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
			},
			// Document Template/materials-col/material-type (PWG5100.21)
			"material-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
		}},
	},
	// Document Template/media (PWG5100.5)
	"media": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Template/media-col (PWG5100.5)
	"media-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Template/media-col/media-top-offset (IPPLABEL)
			"media-top-offset": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   -2147483648,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/media-col/media-tracking (IPPLABEL)
			"media-tracking": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// Document Template/media-input-tray-check (PWG5100.5)
	"media-input-tray-check": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Template/media-overprint (PWG5100.13)
	"media-overprint": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Template/media-overprint/media-overprint-distance (PWG5100.13)
			"media-overprint-distance": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/media-overprint/media-overprint-method (PWG5100.13)
			"media-overprint-method": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// Document Template/multiple-object-handling (PWG5100.21)
	"multiple-object-handling": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/number-up (PWG5100.5)
	"number-up": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/orientation-requested (PWG5100.5)
	"orientation-requested": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Template/output-bin (PWG5100.5)
	"output-bin": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Document Template/output-device (PWG5100.7)
	"output-device": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Document Template/page-delivery (PWG5100.5)
	"page-delivery": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/page-order-received (IPP20190509B)
	"page-order-received": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/page-ranges (PWG5100.5)
	"page-ranges": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Document Template/platform-temperature (PWG5100.21)
	"platform-temperature": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/presentation-direction-number-up (PWG5100.5)
	"presentation-direction-number-up": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-accuracy (PWG5100.21)
	"print-accuracy": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Template/print-accuracy/accuracy-units (PWG5100.21)
			"accuracy-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Document Template/print-accuracy/x-accuracy (PWG5100.21)
			"x-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/print-accuracy/y-accuracy (PWG5100.21)
			"y-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/print-accuracy/z-accuracy (PWG5100.21)
			"z-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Document Template/print-base (PWG5100.21)
	"print-base": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-color-mode (PWG5100.13)
	"print-color-mode": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-content-optimize (PWG5100.7)
	"print-content-optimize": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-darkness (IPPLABEL)
	"print-darkness": &DefAttr{
		SetOf: false,
		Min:   -100,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/print-objects (PWG5100.21)
	"print-objects": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Document Template/print-objects/document-number (PWG5100.21)
			"document-number": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Document Template/print-objects/object-offset (PWG5100.21)
			"object-offset": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Document Template/print-objects/object-offset/x-offset (PWG5100.21)
					"x-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Document Template/print-objects/object-offset/y-offset (PWG5100.21)
					"y-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Document Template/print-objects/object-offset/z-offset (PWG5100.21)
					"z-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Document Template/print-objects/object-size (PWG5100.21)
			"object-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Document Template/print-objects/object-size/x-dimension (PWG5100.21)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Document Template/print-objects/object-size/y-dimension (PWG5100.21)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Document Template/print-objects/object-size/z-dimension (PWG5100.21)
					"z-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Document Template/print-objects/object-uuid (PWG5100.21)
			"object-uuid": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Document Template/print-quality (PWG5100.5)
	"print-quality": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Document Template/print-rendering-intent (PWG5100.13)
	"print-rendering-intent": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-scaling (PWG5100.13)
	"print-scaling": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/print-speed (IPPLABEL)
	"print-speed": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/print-supports (PWG5100.21)
	"print-supports": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/printer-resolution (PWG5100.5)
	"printer-resolution": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Document Template/separator-sheets (PWG5100.5)
	"separator-sheets": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Document Template/sheet-collate (PWG5100.5)
	"sheet-collate": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/sides (PWG5100.5)
	"sides": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/x-image-position (PWG5100.5)
	"x-image-position": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/x-image-shift (PWG5100.5)
	"x-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/x-side1-image-shift (PWG5100.5)
	"x-side1-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/x-side2-image-shift (PWG5100.5)
	"x-side2-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/y-image-position (PWG5100.5)
	"y-image-position": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Document Template/y-image-shift (PWG5100.5)
	"y-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/y-side1-image-shift (PWG5100.5)
	"y-side1-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Document Template/y-side2-image-shift (PWG5100.5)
	"y-side2-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
}

// EventNotifications is the Event Notifications attributes
var EventNotifications = map[string]*DefAttr{
	// Event Notifications/job-id (rfc3996)
	"job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Event Notifications/job-impressions-completed (rfc3996)
	"job-impressions-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Event Notifications/job-state (rfc3996)
	"job-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Event Notifications/job-state-reasons (rfc3996)
	"job-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Event Notifications/job-uuid (PWG5100.13)
	"job-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Event Notifications/notify-charset (rfc3996)
	"notify-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Event Notifications/notify-natural-language (rfc3996)
	"notify-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Event Notifications/notify-printer-uri (rfc3996)
	"notify-printer-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Event Notifications/notify-sequence-number (rfc3996)
	"notify-sequence-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Event Notifications/notify-subscribed-event (rfc3995)
	"notify-subscribed-event": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Event Notifications/notify-subscription-id (rfc3996)
	"notify-subscription-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Event Notifications/notify-subscription-uuid (PWG5100.13)
	"notify-subscription-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Event Notifications/notify-text (rfc3995)
	"notify-text": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Event Notifications/notify-user-data (rfc3996)
	"notify-user-data": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Event Notifications/printer-current-time (rfc3996)
	"printer-current-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagUnknown},
	},
	// Event Notifications/printer-is-accepting-jobs (rfc3996)
	"printer-is-accepting-jobs": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Event Notifications/printer-state (rfc3996)
	"printer-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Event Notifications/printer-state-reasons (rfc3996)
	"printer-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Event Notifications/printer-up-time (rfc3996)
	"printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
}

// JobDescription is the Job Description attributes
var JobDescription = map[string]*DefAttr{
	// Job Description/auth-info (CUPS)
	"auth-info": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Description/current-page-order (IPP20190509B)
	"current-page-order": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Description/job-charge-info (PWG5100.16)
	"job-charge-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Description/job-collation-type (rfc3381)
	"job-collation-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Description/job-impressions-col (PWG5100.7)
	"job-impressions-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Description/job-impressions-col/blank (PWG5100.7)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/blank-two-sided (PWG5100.7)
			"blank-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/full-color-two-sided (PWG5100.7)
			"full-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/highlight-color (PWG5100.7)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/highlight-color-two-sided (PWG5100.7)
			"highlight-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/monochrome (PWG5100.7)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-impressions-col/monochrome-two-sided (PWG5100.7)
			"monochrome-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Description/job-media-progress (CUPS)
	"job-media-progress": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Description/job-media-sheets-col (PWG5100.7)
	"job-media-sheets-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Description/job-media-sheets-col/blank (PWG5100.7)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/blank-two-sided (PWG5100.7)
			"blank-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/full-color-two-sided (PWG5100.7)
			"full-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/highlight-color (PWG5100.7)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/highlight-color-two-sided (PWG5100.7)
			"highlight-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/monochrome (PWG5100.7)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Description/job-media-sheets-col/monochrome-two-sided (PWG5100.7)
			"monochrome-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Description/job-message-from-operator (rfc8011)
	"job-message-from-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Description/job-message-to-operator-actual (PWG5100.8)
	"job-message-to-operator-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Description/job-name (rfc8011)
	"job-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Description/job-originating-host-name (CUPS)
	"job-originating-host-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Description/job-printer-state-message (CUPS)
	"job-printer-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Description/job-printer-state-reasons (CUPS)
	"job-printer-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Description/job-save-printer-make-and-model (PWG5100.11)
	"job-save-printer-make-and-model": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
}

// JobStatus is the Job Status attributes
var JobStatus = map[string]*DefAttr{
	// Job Status/attributes-charset (rfc8011)
	"attributes-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Job Status/attributes-natural-language (rfc8011)
	"attributes-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Job Status/chamber-humidity-actual (PWG5100.21)
	"chamber-humidity-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/chamber-temperature-actual (PWG5100.21)
	"chamber-temperature-actual": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/client-info (PWG5100.7)
	"client-info": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/compression-supplied (PWG5100.7)
	"compression-supplied": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/copies-actual (PWG5100.8)
	"copies-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/cover-back-actual (PWG5100.8)
	"cover-back-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/cover-front-actual (PWG5100.8)
	"cover-front-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/current-page-order (PWG5100.3)
	"current-page-order": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/date-time-at-completed (rfc8011)
	"date-time-at-completed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagNoValue},
	},
	// Job Status/date-time-at-completed-estimated (PWG5100.3)
	"date-time-at-completed-estimated": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagUnknown},
	},
	// Job Status/date-time-at-creation (rfc8011)
	"date-time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Job Status/date-time-at-processing (rfc8011)
	"date-time-at-processing": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagNoValue},
	},
	// Job Status/date-time-at-processing-estimated (PWG5100.3)
	"date-time-at-processing-estimated": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagUnknown},
	},
	// Job Status/destination-statuses (PWG5100.15)
	"destination-statuses": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Status/destination-statuses/destination-uri (PWG5100.15)
			"destination-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Job Status/destination-statuses/images-completed (PWG5100.15)
			"images-completed": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/destination-statuses/transmission-status (PWG5100.15)
			"transmission-status": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
		}},
	},
	// Job Status/document-charset-supplied (PWG5100.7)
	"document-charset-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Job Status/document-digital-signature-supplied (PWG5100.7)
	"document-digital-signature-supplied": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/document-format-details-detected (PWG5100.7)
	"document-format-details-detected": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/document-format-details-supplied (PWG5100.7-2003)
	"document-format-details-supplied": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/document-format-ready (PWG5100.18)
	"document-format-ready": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Job Status/document-format-supplied (PWG5100.7)
	"document-format-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Job Status/document-format-version-supplied (PWG5100.7)
	"document-format-version-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/document-message-supplied (PWG5100.7)
	"document-message-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/document-metadata (PWG5100.13)
	"document-metadata": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Job Status/document-name-supplied (PWG5100.7)
	"document-name-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/document-natural-language-supplied (PWG5100.7)
	"document-natural-language-supplied": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Job Status/errors-count (PWG5100.7)
	"errors-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/finishings-actual (PWG5100.8)
	"finishings-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Status/finishings-col-actual (PWG5100.8)
	"finishings-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/force-front-side-actual (PWG5100.8)
	"force-front-side-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/imposition-template-actual (PWG5100.8)
	"imposition-template-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/impressions-completed-current-copy (rfc3381)
	"impressions-completed-current-copy": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/input-attributes-actual (PWG5100.15)
	"input-attributes-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/insert-sheet-actual (PWG5100.8)
	"insert-sheet-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/ipp-attribute-fidelity (PWG5100.7)
	"ipp-attribute-fidelity": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Job Status/job-account-id-actual (PWG5100.8)
	"job-account-id-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/job-account-type-actual (PWG5100.16)
	"job-account-type-actual": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/job-accounting-sheets-actual (PWG5100.8)
	"job-accounting-sheets-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-accounting-user-id-actual (PWG5100.8)
	"job-accounting-user-id-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/job-copies-actual (PWG5100.7)
	"job-copies-actual": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-cover-back-actual (PWG5100.7)
	"job-cover-back-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-cover-front-actual (PWG5100.7)
	"job-cover-front-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-detailed-status-messages (rfc8011)
	"job-detailed-status-messages": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/job-document-access-errors (rfc8011)
	"job-document-access-errors": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/job-error-sheet-actual (PWG5100.8)
	"job-error-sheet-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-finishings-actual (PWG5100.7)
	"job-finishings-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Status/job-hold-until-actual (PWG5100.8)
	"job-hold-until-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/job-id (rfc8011)
	"job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-impressions (rfc8011)
	"job-impressions": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-impressions-col (XEROX20150505)
	"job-impressions-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Status/job-impressions-col/blank (XEROX20150505)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/blank-two-sided (XEROX20150505)
			"blank-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/full-color (XEROX20150505)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/full-color-two-sided (XEROX20150505)
			"full-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/highlight-color (XEROX20150505)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/highlight-color-two-sided (XEROX20150505)
			"highlight-color-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/monochrome (XEROX20150505)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-impressions-col/monochrome-two-sided (XEROX20150505)
			"monochrome-two-sided": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Status/job-impressions-completed (rfc8011)
	"job-impressions-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-impressions-completed-col (PWG5100.7)
	"job-impressions-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-k-octets (rfc8011)
	"job-k-octets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-k-octets-processed (rfc8011)
	"job-k-octets-processed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-mandatory-attributes (PWG5100.7)
	"job-mandatory-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/job-media-sheets (rfc8011)
	"job-media-sheets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-media-sheets-col (XEROX20150505)
	"job-media-sheets-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Status/job-media-sheets-col/blank (XEROX20150505)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-media-sheets-col/full-color (XEROX20150505)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-media-sheets-col/highlight-color (XEROX20150505)
			"highlight-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-media-sheets-col/monochrome (XEROX20150505)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Status/job-media-sheets-completed (rfc8011)
	"job-media-sheets-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-media-sheets-completed-col (PWG5100.7)
	"job-media-sheets-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-more-info (rfc8011)
	"job-more-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/job-originating-user-name (rfc8011)
	"job-originating-user-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/job-originating-user-uri (PWG5100.13)
	"job-originating-user-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/job-pages (PWG5100.13)
	"job-pages": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-pages-col (PWG5100.7)
	"job-pages-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Status/job-pages-col/blank (PWG5100.7)
			"blank": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-pages-col/full-color (PWG5100.7)
			"full-color": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Status/job-pages-col/monochrome (PWG5100.7)
			"monochrome": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Status/job-pages-completed (PWG5100.13)
	"job-pages-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-pages-completed-col (PWG5100.7)
	"job-pages-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-printer-up-time (rfc8011)
	"job-printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-printer-uri (rfc8011)
	"job-printer-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/job-priority-actual (PWG5100.8)
	"job-priority-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-processing-time (PWG5100.7)
	"job-processing-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-release-action (PWG5100.11)
	"job-release-action": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/job-resource-ids (PWG5100.22)
	"job-resource-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/job-sheet-message-actual (PWG5100.8)
	"job-sheet-message-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/job-sheets-actual (PWG5100.8)
	"job-sheets-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/job-sheets-col-actual (PWG5100.8)
	"job-sheets-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-state (rfc8011)
	"job-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum, goipp.TagUnknown},
	},
	// Job Status/job-state-message (rfc8011)
	"job-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/job-state-reasons (rfc8011)
	"job-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/job-storage (PWG5100.11)
	"job-storage": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/job-uri (rfc8011)
	"job-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/job-uuid (PWG5100.13)
	"job-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/materials-col-actual (PWG5100.21)
	"materials-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/media-actual (PWG5100.8)
	"media-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/media-col-actual (PWG5100.8)
	"media-col-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/media-input-tray-check-actual (PWG5100.8)
	"media-input-tray-check-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/multiple-document-handling-actual (PWG5100.8)
	"multiple-document-handling-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/multiple-object-handling-actual (PWG5100.21)
	"multiple-object-handling-actual": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/number-of-documents (rfc8011)
	"number-of-documents": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/number-of-intervening-jobs (rfc8011)
	"number-of-intervening-jobs": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/number-up-actual (PWG5100.8)
	"number-up-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/orientation-requested-actual (PWG5100.8)
	"orientation-requested-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Status/original-requesting-user-name (rfc3998)
	"original-requesting-user-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/output-attributes-actual (PWG5100.17)
	"output-attributes-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/output-bin-actual (PWG5100.8)
	"output-bin-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Status/output-device-actual (PWG5100.7)
	"output-device-actual": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/output-device-assigned (rfc8011)
	"output-device-assigned": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Status/output-device-job-state (PWG5100.18)
	"output-device-job-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Status/output-device-job-state-message (PWG5100.18)
	"output-device-job-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Status/output-device-job-state-reasons (PWG5100.18)
	"output-device-job-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/output-device-uuid-assigned (PWG5100.18)
	"output-device-uuid-assigned": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/overrides-actual (PWG5100.6)
	"overrides-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/page-delivery-actual (PWG5100.8)
	"page-delivery-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/page-order-received-actual (IPP20190509B)
	"page-order-received-actual": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/page-ranges-actual (PWG5100.8)
	"page-ranges-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Job Status/parent-job-id (PWG5100.11)
	"parent-job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/parent-job-uuid (PWG5100.11)
	"parent-job-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   45,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Status/platform-temperature-actual (PWG5100.21)
	"platform-temperature-actual": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/presentation-direction-number-up-actual (PWG5100.8)
	"presentation-direction-number-up-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/print-accuracy-actual (PWG5100.21)
	"print-accuracy-actual": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/print-base-actual (PWG5100.21)
	"print-base-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/print-color-mode-actual (PWG5100.13)
	"print-color-mode-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/print-content-optimize-actual (PWG5100.7)
	"print-content-optimize-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/print-objects-actual (PWG5100.21)
	"print-objects-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/print-quality-actual (PWG5100.8)
	"print-quality-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Status/print-supports-actual (PWG5100.21)
	"print-supports-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/printer-resolution-actual (PWG5100.8)
	"printer-resolution-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Job Status/separator-sheets-actual (PWG5100.8)
	"separator-sheets-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Status/sheet-collate-actual (PWG5100.8)
	"sheet-collate-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/sheet-completed-copy-number (rfc3381)
	"sheet-completed-copy-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/sheet-completed-document-number (rfc3381)
	"sheet-completed-document-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/sides-actual (PWG5100.8)
	"sides-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/time-at-completed (rfc8011)
	"time-at-completed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Job Status/time-at-completed-estimated (PWG5100.3)
	"time-at-completed-estimated": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagUnknown},
	},
	// Job Status/time-at-creation (rfc8011)
	"time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/time-at-processing (rfc8011)
	"time-at-processing": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Job Status/time-at-processing-estimated (PWG5100.3)
	"time-at-processing-estimated": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagUnknown},
	},
	// Job Status/warnings-count (PWG5100.7)
	"warnings-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/x-image-position-actual (PWG5100.8)
	"x-image-position-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/x-image-shift-actual (PWG5100.8)
	"x-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/x-side1-image-shift-actual (PWG5100.8)
	"x-side1-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/x-side2-image-shift-actual (PWG5100.8)
	"x-side2-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/y-image-position-actual (PWG5100.8)
	"y-image-position-actual": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Status/y-image-shift-actual (PWG5100.8)
	"y-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/y-side1-image-shift-actual (PWG5100.8)
	"y-side1-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Status/y-side2-image-shift-actual (PWG5100.8)
	"y-side2-image-shift-actual": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
}

// JobTemplate is the Job Template attributes
var JobTemplate = map[string]*DefAttr{
	// Job Template/chamber-humidity (PWG5100.21)
	"chamber-humidity": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/chamber-temperature (PWG5100.21)
	"chamber-temperature": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/confirmation-sheet-print (PWG5100.15)
	"confirmation-sheet-print": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Job Template/copies (rfc8011)
	"copies": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/cover-back (PWG5100.3)
	"cover-back": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/cover-back/cover-type (PWG5100.3)
			"cover-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/cover-back/media (rfc8011)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/cover-back/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/cover-front (PWG5100.3)
	"cover-front": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/cover-front/cover-type (PWG5100.3)
			"cover-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/cover-front/media (rfc8011)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/cover-front/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/cover-sheet-info (PWG5100.15)
	"cover-sheet-info": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/cover-sheet-info/from-name (PWG5100.15)
			"from-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/cover-sheet-info/logo (PWG5100.15)
			"logo": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Job Template/cover-sheet-info/message (PWG5100.15)
			"message": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/cover-sheet-info/organization-name (PWG5100.15)
			"organization-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/cover-sheet-info/subject (PWG5100.15)
			"subject": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/cover-sheet-info/to-name (PWG5100.15)
			"to-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
		}},
	},
	// Job Template/destination-uris (PWG5100.15)
	"destination-uris": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/destination-uris/destination-attributes (PWG5100.17)
			"destination-attributes": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
			// Job Template/destination-uris/destination-uri (PWG5100.15)
			"destination-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Job Template/destination-uris/post-dial-string (PWG5100.15)
			"post-dial-string": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/destination-uris/pre-dial-string (PWG5100.15)
			"pre-dial-string": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/destination-uris/t33-subaddress (PWG5100.15)
			"t33-subaddress": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Template/feed-orientation (PWG5100.11)
	"feed-orientation": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/finishings (rfc8011)
	"finishings": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Template/finishings-col (PWG5100.1)
	"finishings-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
		Members: []map[string]*DefAttr{{
			// Job Template/finishings-col/baling (PWG5100.1)
			"baling": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/baling/baling-type (PWG5100.1)
					"baling-type": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
					// Job Template/finishings-col/baling/baling-when (PWG5100.1)
					"baling-when": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
				}},
			},
			// Job Template/finishings-col/binding (PWG5100.1)
			"binding": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/binding/binding-reference-edge (PWG5100.1)
					"binding-reference-edge": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/binding/binding-type (PWG5100.1)
					"binding-type": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
				}},
			},
			// Job Template/finishings-col/coating (PWG5100.1)
			"coating": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/coating/coating-sides (PWG5100.1)
					"coating-sides": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/coating/coating-type (PWG5100.1)
					"coating-type": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
				}},
			},
			// Job Template/finishings-col/covering (PWG5100.1)
			"covering": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/covering/covering-name (PWG5100.1)
					"covering-name": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
				}},
			},
			// Job Template/finishings-col/finishing-template (PWG5100.1)
			"finishing-template": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/finishings-col/folding (PWG5100.1)
			"folding": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/folding/folding-direction (PWG5100.1)
					"folding-direction": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/folding/folding-offset (PWG5100.1)
					"folding-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/folding/folding-reference-edge (PWG5100.1)
					"folding-reference-edge": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
				}},
			},
			// Job Template/finishings-col/imposition-template (PWG5100.1)
			"imposition-template": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/finishings-col/laminating (PWG5100.1)
			"laminating": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/laminating/laminating-sides (PWG5100.1)
					"laminating-sides": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/laminating/laminating-type (PWG5100.1)
					"laminating-type": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
				}},
			},
			// Job Template/finishings-col/media-sheets-supported (PWG5100.1)
			"media-sheets-supported": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Job Template/finishings-col/media-size (PWG5100.1)
			"media-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
			// Job Template/finishings-col/media-size-name (PWG5100.1)
			"media-size-name": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/finishings-col/punching (PWG5100.1)
			"punching": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/punching/punching-locations (PWG5100.1)
					"punching-locations": &DefAttr{
						SetOf: true,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/punching/punching-offset (PWG5100.1)
					"punching-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/punching/punching-reference-edge (PWG5100.1)
					"punching-reference-edge": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
				}},
			},
			// Job Template/finishings-col/stitching (PWG5100.1)
			"stitching": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/stitching/stitching-angle (PWG5100.1)
					"stitching-angle": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   359,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/stitching/stitching-locations (PWG5100.1)
					"stitching-locations": &DefAttr{
						SetOf: true,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/stitching/stitching-method (PWG5100.1)
					"stitching-method": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/stitching/stitching-offset (PWG5100.1)
					"stitching-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/stitching/stitching-reference-edge (PWG5100.1)
					"stitching-reference-edge": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
				}},
			},
			// Job Template/finishings-col/trimming (PWG5100.1)
			"trimming": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/finishings-col/trimming/trimming-offset (PWG5100.1)
					"trimming-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/finishings-col/trimming/trimming-reference-edge (PWG5100.1)
					"trimming-reference-edge": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Job Template/finishings-col/trimming/trimming-type (PWG5100.1)
					"trimming-type": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
					},
					// Job Template/finishings-col/trimming/trimming-when (PWG5100.1)
					"trimming-when": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
				}},
			},
		}},
	},
	// Job Template/force-front-side (PWG5100.3)
	"force-front-side": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/image-orientation (PWG5100.3)
	"image-orientation": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Template/imposition-template (PWG5100.3)
	"imposition-template": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/insert-sheet (PWG5100.3)
	"insert-sheet": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/insert-sheet/insert-after-page-number (PWG5100.3)
			"insert-after-page-number": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/insert-sheet/insert-count (PWG5100.3)
			"insert-count": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/insert-sheet/media (PWG5100.3)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/insert-sheet/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/job-account-id (PWG5100.7)
	"job-account-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Template/job-account-type (PWG5100.16)
	"job-account-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-accounting-sheets (PWG5100.3)
	"job-accounting-sheets": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/job-accounting-sheets/job-accounting-sheets-type (PWG5100.3)
			"job-accounting-sheets-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-accounting-sheets/media (PWG5100.3)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-accounting-sheets/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/job-accounting-user-id (PWG5100.7)
	"job-accounting-user-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Template/job-cancel-after (PWG5100.11)
	"job-cancel-after": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/job-complete-before (PWG5100.3)
	"job-complete-before": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-complete-before-time (PWG5100.3)
	"job-complete-before-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Job Template/job-copies (PWG5100.7)
	"job-copies": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/job-cover-back (PWG5100.7)
	"job-cover-back": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Template/job-cover-front (PWG5100.7)
	"job-cover-front": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Job Template/job-delay-output-until (PWG5100.7)
	"job-delay-output-until": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-delay-output-until-time (PWG5100.7)
	"job-delay-output-until-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Job Template/job-error-action (PWG5100.13)
	"job-error-action": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/job-error-sheet (PWG5100.3)
	"job-error-sheet": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/job-error-sheet/job-error-sheet-type (PWG5100.3)
			"job-error-sheet-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-error-sheet/job-error-sheet-when (PWG5100.3)
			"job-error-sheet-when": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/job-error-sheet/media (PWG5100.3)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-error-sheet/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/job-finishings (PWG5100.7)
	"job-finishings": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Template/job-hold-until (rfc8011)
	"job-hold-until": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-hold-until-time (PWG5100.7)
	"job-hold-until-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Job Template/job-message-to-operator (PWG5100.3)
	"job-message-to-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Template/job-pages-per-set (PWG5100.1)
	"job-pages-per-set": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/job-phone-number (PWG5100.3)
	"job-phone-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Job Template/job-priority (rfc8011)
	"job-priority": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/job-recipient-name (PWG5100.3)
	"job-recipient-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Template/job-retain-until (PWG5100.7)
	"job-retain-until": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-retain-until-interval (PWG5100.7)
	"job-retain-until-interval": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/job-retain-until-time (PWG5100.7)
	"job-retain-until-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Job Template/job-sheet-message (PWG5100.3)
	"job-sheet-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Job Template/job-sheets (rfc8011)
	"job-sheets": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/job-sheets-col (PWG5100.7)
	"job-sheets-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/job-sheets-col/job-sheets (PWG5100.7)
			"job-sheets": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-sheets-col/media (PWG5100.7)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/job-sheets-col/media-col (PWG5100.7)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// Job Template/materials-col (PWG5100.21)
	"materials-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/materials-col/material-amount (PWG5100.21)
			"material-amount": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-amount-units (PWG5100.21)
			"material-amount-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/materials-col/material-color (PWG5100.21)
			"material-color": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/materials-col/material-diameter (PWG5100.21)
			"material-diameter": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-diameter-tolerance (PWG5100.21)
			"material-diameter-tolerance": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-fill-density (PWG5100.21)
			"material-fill-density": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-key (PWG5100.21)
			"material-key": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/materials-col/material-name (PWG5100.21)
			"material-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Job Template/materials-col/material-nozzle-diameter (PWG5100.21)
			"material-nozzle-diameter": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-purpose (PWG5100.21)
			"material-purpose": &DefAttr{
				SetOf: true,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/materials-col/material-rate (PWG5100.21)
			"material-rate": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-rate-units (PWG5100.21)
			"material-rate-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/materials-col/material-retraction (PWG5100.21)
			"material-retraction": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Job Template/materials-col/material-shell-thickness (PWG5100.21)
			"material-shell-thickness": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/materials-col/material-temperature (PWG5100.21)
			"material-temperature": &DefAttr{
				SetOf: false,
				Min:   -273,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
			},
			// Job Template/materials-col/material-type (PWG5100.21)
			"material-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
		}},
	},
	// Job Template/media (rfc8011)
	"media": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/media-col (PWG5100.7)
	"media-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/media-col/media-back-coating (PWG5100.7)
			"media-back-coating": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-bottom-margin (PWG5100.7)
			"media-bottom-margin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-color (PWG5100.7)
			"media-color": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-front-coating (PWG5100.7)
			"media-front-coating": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-grain (PWG5100.7)
			"media-grain": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-hole-count (PWG5100.7)
			"media-hole-count": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-info (PWG5100.7)
			"media-info": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Job Template/media-col/media-key (PWG5100.7)
			"media-key": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-left-margin (PWG5100.7)
			"media-left-margin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-order-count (PWG5100.7)
			"media-order-count": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-pre-printed (PWG5100.7)
			"media-pre-printed": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-recycled (PWG5100.7)
			"media-recycled": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-right-margin (PWG5100.7)
			"media-right-margin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-size (PWG5100.7)
			"media-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/media-col/media-size/x-dimension (PWG5100.7)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/media-col/media-size/y-dimension (PWG5100.7)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Job Template/media-col/media-size-name (PWG5100.7)
			"media-size-name": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-source (PWG5100.7)
			"media-source": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-thickness (PWG5100.7)
			"media-thickness": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-tooth (PWG5100.7)
			"media-tooth": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-top-margin (PWG5100.7)
			"media-top-margin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-top-offset (IPPLABEL)
			"media-top-offset": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   -2147483648,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-col/media-tracking (IPPLABEL)
			"media-tracking": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/media-col/media-type (PWG5100.7)
			"media-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/media-col/media-weight-metric (PWG5100.7)
			"media-weight-metric": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Template/media-input-tray-check (PWG5100.3)
	"media-input-tray-check": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/media-overprint (PWG5100.13)
	"media-overprint": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/media-overprint/media-overprint-distance (PWG5100.13)
			"media-overprint-distance": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/media-overprint/media-overprint-method (PWG5100.13)
			"media-overprint-method": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// Job Template/multiple-document-handling (rfc8011)
	"multiple-document-handling": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/multiple-object-handling (PWG5100.21)
	"multiple-object-handling": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/number-of-retries (PWG5100.15)
	"number-of-retries": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/number-up (rfc8011)
	"number-up": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/orientation-requested (rfc8011)
	"orientation-requested": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Template/output-bin (PWG5100.2)
	"output-bin": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Job Template/output-device (PWG5100.7)
	"output-device": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Job Template/overrides (PWG5100.6)
	"overrides": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/overrides/document-copies (PWG5100.6)
			"document-copies": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Job Template/overrides/document-numbers (PWG5100.6)
			"document-numbers": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Job Template/overrides/pages (PWG5100.6)
			"pages": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
		}},
	},
	// Job Template/page-border (CUPS)
	"page-border": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/page-delivery (PWG5100.3)
	"page-delivery": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/page-order-received (PWG5100.3)
	"page-order-received": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/page-ranges (rfc8011)
	"page-ranges": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Job Template/page-set (CUPS)
	"page-set": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/pages-per-subset (PWG5100.13)
	"pages-per-subset": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/pclm-source-resolution (HP20180907)
	"pclm-source-resolution": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Job Template/platform-temperature (PWG5100.21)
	"platform-temperature": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/presentation-direction-number-up (PWG5100.3)
	"presentation-direction-number-up": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-accuracy (PWG5100.21)
	"print-accuracy": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/print-accuracy/accuracy-units (PWG5100.21)
			"accuracy-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Job Template/print-accuracy/x-accuracy (PWG5100.21)
			"x-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/print-accuracy/y-accuracy (PWG5100.21)
			"y-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/print-accuracy/z-accuracy (PWG5100.21)
			"z-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Template/print-base (PWG5100.21)
	"print-base": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-color-mode (PWG5100.13)
	"print-color-mode": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-content-optimize (PWG5100.7)
	"print-content-optimize": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-darkness (IPPLABEL)
	"print-darkness": &DefAttr{
		SetOf: false,
		Min:   -100,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/print-objects (PWG5100.21)
	"print-objects": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/print-objects/document-number (PWG5100.21)
			"document-number": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Job Template/print-objects/object-offset (PWG5100.21)
			"object-offset": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/print-objects/object-offset/x-offset (PWG5100.21)
					"x-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/print-objects/object-offset/y-offset (PWG5100.21)
					"y-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/print-objects/object-offset/z-offset (PWG5100.21)
					"z-offset": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Job Template/print-objects/object-size (PWG5100.21)
			"object-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Job Template/print-objects/object-size/x-dimension (PWG5100.21)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/print-objects/object-size/y-dimension (PWG5100.21)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Job Template/print-objects/object-size/z-dimension (PWG5100.21)
					"z-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Job Template/print-objects/object-uuid (PWG5100.21)
			"object-uuid": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Job Template/print-quality (rfc8011)
	"print-quality": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Job Template/print-rendering-intent (PWG5100.13)
	"print-rendering-intent": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-scaling (PWG5100.13)
	"print-scaling": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/print-speed (IPPLABEL)
	"print-speed": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/print-supports (PWG5100.21)
	"print-supports": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/printer-resolution (rfc8011)
	"printer-resolution": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Job Template/proof-copies (PWG5100.11)
	"proof-copies": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/proof-print (PWG5100.11)
	"proof-print": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/proof-print/media (PWG5100.11)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/proof-print/media-col (PWG5100.11)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
			// Job Template/proof-print/proof-print-copies (PWG5100.11)
			"proof-print-copies": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Job Template/retry-interval (PWG5100.15)
	"retry-interval": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/retry-time-out (PWG5100.15)
	"retry-time-out": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/separator-sheets (PWG5100.3)
	"separator-sheets": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Job Template/separator-sheets/media (rfc8011)
			"media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Job Template/separator-sheets/media-col (PWG5100.3)
			"media-col": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
			// Job Template/separator-sheets/separator-sheets-type (PWG5100.3)
			"separator-sheets-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
		}},
	},
	// Job Template/sheet-collate (rfc3381)
	"sheet-collate": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/sides (rfc8011)
	"sides": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/x-image-position (PWG5100.3)
	"x-image-position": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/x-image-shift (PWG5100.3)
	"x-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/x-side1-image-shift (PWG5100.3)
	"x-side1-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/x-side2-image-shift (PWG5100.3)
	"x-side2-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/y-image-position (PWG5100.3)
	"y-image-position": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Job Template/y-image-shift (PWG5100.3)
	"y-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/y-side1-image-shift (PWG5100.3)
	"y-side1-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Job Template/y-side2-image-shift (PWG5100.3)
	"y-side2-image-shift": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
}

// Operation is the Operation attributes
var Operation = map[string]*DefAttr{
	// Operation/attributes-charset (rfc8011)
	"attributes-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Operation/attributes-natural-language (rfc8011)
	"attributes-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Operation/charge-info-message (PWG5100.16)
	"charge-info-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/client-info (PWG5100.7)
	"client-info": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/client-info/client-name (PWG5100.7)
			"client-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Operation/client-info/client-patches (PWG5100.7)
			"client-patches": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagText, goipp.TagNoValue},
			},
			// Operation/client-info/client-string-version (PWG5100.7)
			"client-string-version": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/client-info/client-type (PWG5100.7)
			"client-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
			// Operation/client-info/client-version (PWG5100.7)
			"client-version": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   64,
				Tags:  []goipp.Tag{goipp.TagString, goipp.TagNoValue},
			},
		}},
	},
	// Operation/compression (rfc8011)
	"compression": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/compression-accepted (PWG5100.17)
	"compression-accepted": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/destination-accesses (PWG5100.17)
	"destination-accesses": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
		Members: []map[string]*DefAttr{{
			// Operation/destination-accesses/access-oauth-token (PWG5100.17)
			"access-oauth-token": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
			// Operation/destination-accesses/access-oauth-uri (PWG5100.17)
			"access-oauth-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Operation/destination-accesses/access-password (PWG5100.17)
			"access-password": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/destination-accesses/access-pin (PWG5100.17)
			"access-pin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/destination-accesses/access-user-name (PWG5100.17)
			"access-user-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/destination-accesses/access-x509-certificate (IPPWG20180620)
			"access-x509-certificate": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
		}},
	},
	// Operation/detailed-status-message (rfc8011)
	"detailed-status-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/device-class (CUPS)
	"device-class": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/document-access (PWG5100.18)
	"document-access": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
		Members: []map[string]*DefAttr{{
			// Operation/document-access/access-oauth-token (PWG5100.18)
			"access-oauth-token": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
			// Operation/document-access/access-oauth-uri (PWG5100.18)
			"access-oauth-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Operation/document-access/access-password (PWG5100.18)
			"access-password": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-access/access-pin (PWG5100.18)
			"access-pin": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-access/access-user-name (PWG5100.18)
			"access-user-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-access/access-x509-certificate (IPPWG20180620)
			"access-x509-certificate": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
		}},
	},
	// Operation/document-access-error (rfc8011)
	"document-access-error": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/document-charset (PWG5100.5)
	"document-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Operation/document-data-get-interval (PWG5100.17)
	"document-data-get-interval": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/document-data-wait (PWG5100.17)
	"document-data-wait": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/document-digital-signature (PWG5100.7)
	"document-digital-signature": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/document-format (rfc8011)
	"document-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/document-format-accepted (PWG5100.18)
	"document-format-accepted": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/document-format-details (PWG5100.7-2003)
	"document-format-details": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/document-format-details/document-format (PWG5100.7-2003)
			"document-format": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagMimeType},
			},
			// Operation/document-format-details/document-format-device-id (PWG5100.7-2003)
			"document-format-device-id": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-format-details/document-format-version (PWG5100.7-2003)
			"document-format-version": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-format-details/document-natural-language (PWG5100.7-2003)
			"document-natural-language": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   63,
				Tags:  []goipp.Tag{goipp.TagLanguage},
			},
			// Operation/document-format-details/document-source-application-name (PWG5100.7-2003)
			"document-source-application-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Operation/document-format-details/document-source-application-version (PWG5100.7-2003)
			"document-source-application-version": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Operation/document-format-details/document-source-os-name (PWG5100.7-2003)
			"document-source-os-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   40,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Operation/document-format-details/document-source-os-version (PWG5100.7-2003)
			"document-source-os-version": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   40,
				Tags:  []goipp.Tag{goipp.TagText},
			},
		}},
	},
	// Operation/document-format-version (PWG5100.7)
	"document-format-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/document-message (PWG5100.5)
	"document-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/document-metadata (PWG5100.13)
	"document-metadata": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Operation/document-name (rfc8011)
	"document-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/document-natural-language (rfc8011)
	"document-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Operation/document-number (PWG5100.5)
	"document-number": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/document-password (PWG5100.13)
	"document-password": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Operation/document-preprocessed (PWG5100.18)
	"document-preprocessed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/document-uri (rfc8011)
	"document-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/encrypted-job-request-format (PWG5100.TRUSTNOONE)
	"encrypted-job-request-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/encrypted-job-request-id (PWG5100.TRUSTNOONE)
	"encrypted-job-request-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/exclude-schemes (CUPS)
	"exclude-schemes": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/fetch-status-code (PWG5100.18)
	"fetch-status-code": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/fetch-status-message (PWG5100.18)
	"fetch-status-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/first-index (PWG5100.13)
	"first-index": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/first-printer-name (CUPS)
	"first-printer-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/identify-actions (PWG5100.13)
	"identify-actions": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/include-schemes (CUPS)
	"include-schemes": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/input-attributes (PWG5100.15)
	"input-attributes": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/input-attributes/input-auto-exposure (PWG5100.15)
			"input-auto-exposure": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Operation/input-attributes/input-auto-scaling (PWG5100.15)
			"input-auto-scaling": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Operation/input-attributes/input-auto-skew-correction (PWG5100.15)
			"input-auto-skew-correction": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Operation/input-attributes/input-brightness (PWG5100.15)
			"input-brightness": &DefAttr{
				SetOf: false,
				Min:   -100,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-color-mode (PWG5100.15)
			"input-color-mode": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/input-attributes/input-content-type (PWG5100.15)
			"input-content-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/input-attributes/input-contrast (PWG5100.15)
			"input-contrast": &DefAttr{
				SetOf: false,
				Min:   -100,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-film-scan-mode (PWG5100.15)
			"input-film-scan-mode": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/input-attributes/input-images-to-transfer (PWG5100.15)
			"input-images-to-transfer": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-media (PWG5100.15)
			"input-media": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
			// Operation/input-attributes/input-orientation-requested (PWG5100.15)
			"input-orientation-requested": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
			// Operation/input-attributes/input-quality (PWG5100.15)
			"input-quality": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
			// Operation/input-attributes/input-resolution (PWG5100.15)
			"input-resolution": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagResolution},
			},
			// Operation/input-attributes/input-scaling-height (PWG5100.15)
			"input-scaling-height": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   1000,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-scaling-width (PWG5100.15)
			"input-scaling-width": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   1000,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-scan-regions (PWG5100.15)
			"input-scan-regions": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Operation/input-attributes/input-scan-regions/x-dimension (PWG5100.15)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   MIN,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Operation/input-attributes/input-scan-regions/x-origin (PWG5100.15)
					"x-origin": &DefAttr{
						SetOf: false,
						Min:   MIN,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Operation/input-attributes/input-scan-regions/y-dimension (PWG5100.15)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   MIN,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Operation/input-attributes/input-scan-regions/y-origin (PWG5100.15)
					"y-origin": &DefAttr{
						SetOf: false,
						Min:   MIN,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Operation/input-attributes/input-sharpness (PWG5100.15)
			"input-sharpness": &DefAttr{
				SetOf: false,
				Min:   -100,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/input-attributes/input-sides (PWG5100.15)
			"input-sides": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/input-attributes/input-source (PWG5100.15)
			"input-source": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// Operation/ipp-attribute-fidelity (rfc8011)
	"ipp-attribute-fidelity": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/job-authorization-uri (PWG5100.16)
	"job-authorization-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/job-hold-until (rfc8011)
	"job-hold-until": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Operation/job-hold-until-time (PWG5100.7)
	"job-hold-until-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Operation/job-id (rfc8011)
	"job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-ids (PWG5100.7)
	"job-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-impressions (rfc8011)
	"job-impressions": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-impressions-col (XEROX20150505)
	"job-impressions-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Operation/job-impressions-estimated (PWG5100.16)
	"job-impressions-estimated": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-k-octets (rfc8011)
	"job-k-octets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-mandatory-attributes (PWG5100.7)
	"job-mandatory-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/job-media-sheets (rfc8011)
	"job-media-sheets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-media-sheets-col (XEROX20150505)
	"job-media-sheets-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Operation/job-message-from-operator (rfc3380)
	"job-message-from-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/job-name (rfc8011)
	"job-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/job-pages (PWG5100.7)
	"job-pages": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/job-pages-col (PWG5100.7)
	"job-pages-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Operation/job-password (PWG5100.11)
	"job-password": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Operation/job-password-encryption (PWG5100.11)
	"job-password-encryption": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/job-release-action (PWG5100.11)
	"job-release-action": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/job-state (rfc8011)
	"job-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/job-state-message (rfc8011)
	"job-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/job-state-reasons (rfc8011)
	"job-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/job-storage (PWG5100.11)
	"job-storage": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/job-storage/job-storage-access (PWG5100.11)
			"job-storage-access": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/job-storage/job-storage-disposition (PWG5100.11)
			"job-storage-disposition": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/job-storage/job-storage-group (PWG5100.11)
			"job-storage-group": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
		}},
	},
	// Operation/job-uri (rfc8011)
	"job-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/last-document (rfc8011)
	"last-document": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/limit (rfc8011)
	"limit": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/message (rfc8011)
	"message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/my-jobs (rfc8011)
	"my-jobs": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/notify-get-interval (rfc3996)
	"notify-get-interval": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/notify-printer-ids (PWG5100.22)
	"notify-printer-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/notify-resource-id (PWG5100.22)
	"notify-resource-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/notify-sequence-numbers (rfc3996)
	"notify-sequence-numbers": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/notify-subscription-ids (rfc3996)
	"notify-subscription-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/notify-wait (rfc3996)
	"notify-wait": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Operation/original-requesting-user-name (rfc3998)
	"original-requesting-user-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/output-attributes (PWG5100.17)
	"output-attributes": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/output-attributes/noise-removal (PWG5100.17)
			"noise-removal": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Operation/output-attributes/output-compression-quality-factor (PWG5100.17)
			"output-compression-quality-factor": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   100,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Operation/output-device-job-states (PWG5100.18)
	"output-device-job-states": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/output-device-uuid (PWG5100.18)
	"output-device-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/output-device-x509-certificate (PWG5100.22)
	"output-device-x509-certificate": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/output-device-x509-request (PWG5100.22)
	"output-device-x509-request": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/preferred-attributes (PWG5100.13)
	"preferred-attributes": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Operation/printer-geo-location (PWG5100.22)
	"printer-geo-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/printer-id (PWG5100.22)
	"printer-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/printer-ids (PWG5100.22)
	"printer-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/printer-location (PWG5100.22)
	"printer-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/printer-message-from-operator (rfc3380)
	"printer-message-from-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/printer-service-type (PWG5100.22)
	"printer-service-type": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/printer-type (CUPS)
	"printer-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/printer-type-mask (CUPS)
	"printer-type-mask": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/printer-up-time (rfc3996)
	"printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/printer-uri (rfc8011)
	"printer-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/printer-xri-requested (PWG5100.22)
	"printer-xri-requested": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Operation/printer-xri-requested/xri-authentication (PWG5100.22)
			"xri-authentication": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Operation/printer-xri-requested/xri-security (PWG5100.22)
			"xri-security": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// Operation/profile-uri-actual (PWG5100.16)
	"profile-uri-actual": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/requested-attributes (rfc8011)
	"requested-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/requesting-user-name (rfc8011)
	"requesting-user-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Operation/requesting-user-pkcs7-public-key (PWG5100.TRUSTNOONE)
	"requesting-user-pkcs7-public-key": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/requesting-user-uri (PWG5100.13)
	"requesting-user-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/resource-format (PWG5100.22)
	"resource-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/resource-format-accepted (PWG5100.22)
	"resource-format-accepted": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/resource-formats (PWG5100.22)
	"resource-formats": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Operation/resource-id (PWG5100.22)
	"resource-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/resource-ids (PWG5100.22)
	"resource-ids": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/resource-k-octets (PWG5100.22)
	"resource-k-octets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/resource-natural-language (PWG5100.22)
	"resource-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Operation/resource-patches (PWG5100.22)
	"resource-patches": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText, goipp.TagNoValue},
	},
	// Operation/resource-signature (PWG5100.22)
	"resource-signature": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Operation/resource-states (PWG5100.22)
	"resource-states": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Operation/resource-string-version (PWG5100.22)
	"resource-string-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText, goipp.TagNoValue},
	},
	// Operation/resource-type (PWG5100.22)
	"resource-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/resource-types (PWG5100.22)
	"resource-types": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/resource-version (PWG5100.22)
	"resource-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   64,
		Tags:  []goipp.Tag{goipp.TagString, goipp.TagNoValue},
	},
	// Operation/restart-get-interval (PWG5100.22)
	"restart-get-interval": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/status-message (rfc8011)
	"status-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Operation/system-uri (PWG5100.22)
	"system-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Operation/timeout (CUPS)
	"timeout": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Operation/which-jobs (rfc8011)
	"which-jobs": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Operation/which-printers (PWG5100.22)
	"which-printers": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
}

// PrinterDescription is the Printer Description attributes
var PrinterDescription = map[string]*DefAttr{
	// Printer Description/accuracy-units-supported (PWG5100.21)
	"accuracy-units-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/auth-info-required (CUPS)
	"auth-info-required": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/baling-type-supported (PWG5100.1)
	"baling-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/baling-when-supported (PWG5100.1)
	"baling-when-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/binding-reference-edge-supported (PWG5100.1)
	"binding-reference-edge-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/binding-type-supported (PWG5100.1)
	"binding-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/chamber-humidity-default (PWG5100.21)
	"chamber-humidity-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Printer Description/chamber-humidity-supported (PWG5100.21)
	"chamber-humidity-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/chamber-temperature-default (PWG5100.21)
	"chamber-temperature-default": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Printer Description/chamber-temperature-supported (PWG5100.21)
	"chamber-temperature-supported": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/charset-configured (rfc8011)
	"charset-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Printer Description/charset-supported (rfc8011)
	"charset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Printer Description/client-info-supported (PWG5100.7)
	"client-info-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/coating-sides-supported (PWG5100.1)
	"coating-sides-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/coating-type-supported (PWG5100.1)
	"coating-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/color-supported (rfc8011)
	"color-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/compression-supported (rfc8011)
	"compression-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/confirmation-sheet-print-default (PWG5100.15)
	"confirmation-sheet-print-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/copies-default (rfc8011)
	"copies-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/copies-supported (rfc8011)
	"copies-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/cover-back-default (PWG5100.3)
	"cover-back-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/cover-back-supported (PWG5100.3)
	"cover-back-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/cover-front-default (PWG5100.3)
	"cover-front-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/cover-front-supported (PWG5100.3)
	"cover-front-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/cover-sheet-info-default (PWG5100.15)
	"cover-sheet-info-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/cover-sheet-info-supported (PWG5100.15)
	"cover-sheet-info-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/cover-type-supported (PWG5100.3)
	"cover-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/covering-name-supported (PWG5100.1)
	"covering-name-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/destination-accesses-supported (PWG5100.17)
	"destination-accesses-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/destination-uri-ready (PWG5100.17)
	"destination-uri-ready": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/destination-uri-ready/destination-attributes (PWG5100.17)
			"destination-attributes": &DefAttr{
				SetOf: true,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
			// Printer Description/destination-uri-ready/destination-attributes-supported (PWG5100.17)
			"destination-attributes-supported": &DefAttr{
				SetOf: true,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/destination-uri-ready/destination-info (PWG5100.17)
			"destination-info": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// Printer Description/destination-uri-ready/destination-is-directory (PWG5100.17)
			"destination-is-directory": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// Printer Description/destination-uri-ready/destination-mandatory-access-attributes (PWG5100.17)
			"destination-mandatory-access-attributes": &DefAttr{
				SetOf: true,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/destination-uri-ready/destination-name (PWG5100.17)
			"destination-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Printer Description/destination-uri-ready/destination-oauth-scope (PWG5100.17)
			"destination-oauth-scope": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
			// Printer Description/destination-uri-ready/destination-oauth-token (PWG5100.17)
			"destination-oauth-token": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagString},
			},
			// Printer Description/destination-uri-ready/destination-oauth-uri (PWG5100.17)
			"destination-oauth-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Printer Description/destination-uri-ready/destination-uri (PWG5100.17)
			"destination-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Printer Description/destination-uri-schemes-supported (PWG5100.15)
	"destination-uri-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/destination-uris-supported (PWG5100.15)
	"destination-uris-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/device-uri (CUPS)
	"device-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/document-access-supported (PWG5100.18)
	"document-access-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-charset-default (PWG5100.7)
	"document-charset-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Printer Description/document-charset-supported (PWG5100.7)
	"document-charset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Printer Description/document-creation-attributes-supported (PWG5100.5)
	"document-creation-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-digital-signature-default (PWG5100.7)
	"document-digital-signature-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-digital-signature-supported (PWG5100.7)
	"document-digital-signature-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-format-default (rfc8011)
	"document-format-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/document-format-details-default (PWG5100.7)
	"document-format-details-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/document-format-details-supported (PWG5100.7)
	"document-format-details-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-format-supported (rfc8011)
	"document-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/document-format-version-default (PWG5100.7)
	"document-format-version-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/document-format-version-supported (PWG5100.7)
	"document-format-version-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/document-natural-language-default (PWG5100.7)
	"document-natural-language-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Printer Description/document-natural-language-supported (PWG5100.7)
	"document-natural-language-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Printer Description/document-password-supported (PWG5100.13)
	"document-password-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/document-privacy-attributes (IPPPRIVACY10)
	"document-privacy-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/document-privacy-scope (IPPPRIVACY10)
	"document-privacy-scope": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/feed-orientation-default (PWG5100.11)
	"feed-orientation-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/feed-orientation-supported (PWG5100.11)
	"feed-orientation-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/fetch-document-attributes-supported (PWG5100.18)
	"fetch-document-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/finishing-template-supported (PWG5100.1)
	"finishing-template-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/finishings-col-database (PWG5100.1)
	"finishings-col-database": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/finishings-col-default (PWG5100.1)
	"finishings-col-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/finishings-col-ready (PWG5100.1)
	"finishings-col-ready": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/finishings-col-supported (PWG5100.1)
	"finishings-col-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/finishings-default (rfc8011)
	"finishings-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/finishings-ready (PWG5100.1)
	"finishings-ready": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/finishings-supported (rfc8011)
	"finishings-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/folding-direction-supported (PWG5100.1)
	"folding-direction-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/folding-offset-supported (PWG5100.1)
	"folding-offset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/folding-reference-edge-supported (PWG5100.1)
	"folding-reference-edge-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/force-front-side-default  (PWG5100.3)
	"force-front-side-default ": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/force-front-side-supported (PWG5100.3)
	"force-front-side-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/force-front-side-supported  (PWG5100.3)
	"force-front-side-supported ": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/from-name-supported (PWG5100.15)
	"from-name-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/generated-natural-language-supported (rfc8011)
	"generated-natural-language-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Printer Description/identify-actions-default (PWG5100.13)
	"identify-actions-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/identify-actions-supported (PWG5100.13)
	"identify-actions-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/image-orientation-default (PWG5100.3)
	"image-orientation-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/image-orientation-supported (PWG5100.3)
	"image-orientation-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/imposition-template-default (PWG5100.3)
	"imposition-template-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/imposition-template-supported (PWG5100.3)
	"imposition-template-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/input-attributes-default (PWG5100.15)
	"input-attributes-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/input-attributes-supported (PWG5100.15)
	"input-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/input-color-mode-supported (PWG5100.15)
	"input-color-mode-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/input-content-type-supported (PWG5100.15)
	"input-content-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/input-film-scan-mode-supported (PWG5100.15)
	"input-film-scan-mode-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/input-media-supported (PWG5100.15)
	"input-media-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/input-orientation-requested-supported (PWG5100.15)
	"input-orientation-requested-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/input-quality-supported (PWG5100.15)
	"input-quality-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/input-resolution-supported (PWG5100.15)
	"input-resolution-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/input-scan-regions-supported (PWG5100.15)
	"input-scan-regions-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/input-scan-regions-supported/x-dimension (PWG5100.15)
			"x-dimension": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Printer Description/input-scan-regions-supported/x-origin (PWG5100.15)
			"x-origin": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Printer Description/input-scan-regions-supported/y-dimension (PWG5100.15)
			"y-dimension": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
			// Printer Description/input-scan-regions-supported/y-origin (PWG5100.15)
			"y-origin": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagRange},
			},
		}},
	},
	// Printer Description/input-sides-supported (PWG5100.15)
	"input-sides-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/input-source-supported (PWG5100.15)
	"input-source-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/insert-after-page-number-supported (PWG5100.3)
	"insert-after-page-number-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/insert-count-supported (PWG5100.3)
	"insert-count-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/insert-sheet-default (PWG5100.3)
	"insert-sheet-default": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/insert-sheet-supported (PWG5100.3)
	"insert-sheet-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/ipp-features-supported (PWG5100.13)
	"ipp-features-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/ipp-versions-supported (rfc8011)
	"ipp-versions-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/ippget-event-life (rfc3996)
	"ippget-event-life": &DefAttr{
		SetOf: false,
		Min:   15,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-account-id-default (PWG5100.3)
	"job-account-id-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/job-account-id-supported (PWG5100.3)
	"job-account-id-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-account-type-default (PWG5100.16)
	"job-account-type-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-account-type-supported (PWG5100.16)
	"job-account-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-accounting-output-bin-default (PWG5100.3)
	"job-accounting-output-bin-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-accounting-output-bin-supported (PWG5100.3)
	"job-accounting-output-bin-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-accounting-sheets-default (PWG5100.3)
	"job-accounting-sheets-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/job-accounting-sheets-supported (PWG5100.3)
	"job-accounting-sheets-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-accounting-sheets-type-supported (PWG5100.3)
	"job-accounting-sheets-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-accounting-user-id-default (PWG5100.3)
	"job-accounting-user-id-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/job-accounting-user-id-supported (PWG5100.3)
	"job-accounting-user-id-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-authorization-uri-supported (PWG5100.16)
	"job-authorization-uri-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-cancel-after-default (PWG5100.11)
	"job-cancel-after-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-cancel-after-supported (PWG5100.7)
	"job-cancel-after-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-complete-before-supported (PWG5100.3)
	"job-complete-before-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-constraints-supported (PWG5100.13)
	"job-constraints-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/job-constraints-supported/resolver-name (PWG5100.13)
			"resolver-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
		}},
	},
	// Printer Description/job-copies-supported (PWG5100.7)
	"job-copies-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-cover-back-default (PWG5100.7)
	"job-cover-back-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/job-cover-back-supported (PWG5100.7)
	"job-cover-back-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-cover-front-default (PWG5100.7)
	"job-cover-front-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/job-cover-front-supported (PWG5100.7)
	"job-cover-front-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-creation-attributes-supported (PWG5100.7)
	"job-creation-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-delay-output-until-default (PWG5100.7)
	"job-delay-output-until-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-delay-output-until-interval-supported (PWG5100.7)
	"job-delay-output-until-interval-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-delay-output-until-supported (PWG5100.7)
	"job-delay-output-until-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-delay-output-until-time-supported (PWG5100.7)
	"job-delay-output-until-time-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-destination-spooling-supported (PWG5100.17)
	"job-destination-spooling-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-error-action-default (PWG5100.13)
	"job-error-action-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-error-action-supported (PWG5100.13)
	"job-error-action-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-error-sheet-default (PWG5100.3)
	"job-error-sheet-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/job-error-sheet-supported (PWG5100.3)
	"job-error-sheet-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-error-sheet-type-supported (PWG5100.3)
	"job-error-sheet-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-error-sheet-when-supported (PWG5100.3)
	"job-error-sheet-when-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-finishings-col-supported (PWG5100.7)
	"job-finishings-col-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-finishings-default (PWG5100.7)
	"job-finishings-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/job-finishings-ready (PWG5100.7)
	"job-finishings-ready": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/job-finishings-supported (PWG5100.7)
	"job-finishings-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/job-history-attributes-configured (PWG5100.7)
	"job-history-attributes-configured": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-history-attributes-supported (PWG5100.7)
	"job-history-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-history-interval-configured (PWG5100.7)
	"job-history-interval-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-history-interval-supported (PWG5100.7)
	"job-history-interval-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-hold-until-default (rfc8011)
	"job-hold-until-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-hold-until-supported (rfc8011)
	"job-hold-until-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-hold-until-time-supported (PWG5100.7)
	"job-hold-until-time-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-ids-supported (PWG5100.7)
	"job-ids-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-impressions-supported (rfc8011)
	"job-impressions-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-k-limit (CUPS)
	"job-k-limit": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-k-octets-supported (rfc8011)
	"job-k-octets-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-mandatory-attributes-supported (PWG5100.7)
	"job-mandatory-attributes-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-media-sheets-supported (rfc8011)
	"job-media-sheets-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-message-to-operator-default (PWG5100.3)
	"job-message-to-operator-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/job-message-to-operator-supported (PWG5100.3)
	"job-message-to-operator-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-page-limit (CUPS)
	"job-page-limit": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-pages-per-set-supported (PWG5100.1)
	"job-pages-per-set-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-password-encryption-supported (PWG5100.11)
	"job-password-encryption-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-password-length-supported (PWG5100.11)
	"job-password-length-supported": &DefAttr{
		SetOf: false,
		Min:   4,
		Max:   765,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-password-repertoire-configured (PWG5100.11)
	"job-password-repertoire-configured": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-password-repertoire-supported (PWG5100.11)
	"job-password-repertoire-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-password-supported (PWG5100.11)
	"job-password-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-phone-number-default (PWG5100.3)
	"job-phone-number-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// Printer Description/job-phone-number-scheme-supported (PWG5100.3)
	"job-phone-number-scheme-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/job-phone-number-supported (PWG5100.3)
	"job-phone-number-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-presets-supported (PWG5100.13)
	"job-presets-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/job-presets-supported/preset-category (PWG5100.13)
			"preset-category": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/job-presets-supported/preset-name (PWG5100.13)
			"preset-name": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
		}},
	},
	// Printer Description/job-priority-default (rfc8011)
	"job-priority-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-priority-supported (rfc8011)
	"job-priority-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-privacy-attributes (IPPPRIVACY10)
	"job-privacy-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-privacy-scope (IPPPRIVACY10)
	"job-privacy-scope": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-quota-period (CUPS)
	"job-quota-period": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/job-recipient-name-default (PWG5100.3)
	"job-recipient-name-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/job-recipient-name-supported (PWG5100.3)
	"job-recipient-name-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-release-action-default (PWG5100.11)
	"job-release-action-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-release-action-supported (PWG5100.11)
	"job-release-action-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-resolvers-supported (PWG5100.13)
	"job-resolvers-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/job-resolvers-supported/resolver-name (PWG5100.13)
			"resolver-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
		}},
	},
	// Printer Description/job-retain-until-default (PWG5100.7)
	"job-retain-until-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-retain-until-interval-default (PWG5100.7)
	"job-retain-until-interval-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Printer Description/job-retain-until-interval-supported (PWG5100.7)
	"job-retain-until-interval-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-retain-until-supported (PWG5100.7)
	"job-retain-until-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-retain-until-time-supported (PWG5100.7)
	"job-retain-until-time-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/job-sheet-message-default (PWG5100.3)
	"job-sheet-message-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/job-sheet-message-supported (PWG5100.3)
	"job-sheet-message-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/job-sheets-col-default (PWG5100.3)
	"job-sheets-col-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/job-sheets-col-supported (PWG5100.3)
	"job-sheets-col-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-sheets-default (CUPS)
	"job-sheets-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-sheets-supported (rfc8011)
	"job-sheets-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/job-spooling-supported (PWG5100.7)
	"job-spooling-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-storage-access-supported (PWG5100.11)
	"job-storage-access-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-storage-disposition-supported (PWG5100.11)
	"job-storage-disposition-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-storage-group-supported (PWG5100.11)
	"job-storage-group-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/job-storage-supported (PWG5100.11)
	"job-storage-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/job-triggers-supported (PWG5100.13)
	"job-triggers-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/job-triggers-supported/preset-name (PWG5100.13)
			"preset-name": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
			},
		}},
	},
	// Printer Description/jpeg-features-supported (PWG5100.13)
	"jpeg-features-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/jpeg-k-octets-supported (PWG5100.13)
	"jpeg-k-octets-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/jpeg-x-dimension-supported (PWG5100.13)
	"jpeg-x-dimension-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/jpeg-y-dimension-supported (PWG5100.13)
	"jpeg-y-dimension-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/label-mode-configured (IPPLABEL)
	"label-mode-configured": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/label-mode-supported (IPPLABEL)
	"label-mode-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/label-tear-offset-configured (IPPLABEL)
	"label-tear-offset-configured": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/label-tear-offset-supported (IPPLABEL)
	"label-tear-offset-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/laminating-sides-supported (PWG5100.1)
	"laminating-sides-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/laminating-type-supported (PWG5100.1)
	"laminating-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/logo-uri-formats-supported (PWG5100.15)
	"logo-uri-formats-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/logo-uri-schemes-supported (PWG5100.15)
	"logo-uri-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/marker-change-time (CUPS)
	"marker-change-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/marker-colors (CUPS)
	"marker-colors": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/marker-high-levels (CUPS)
	"marker-high-levels": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/marker-levels (CUPS)
	"marker-levels": &DefAttr{
		SetOf: true,
		Min:   -3,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/marker-low-levels (CUPS)
	"marker-low-levels": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/marker-message (CUPS)
	"marker-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/marker-names (CUPS)
	"marker-names": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/marker-types (CUPS)
	"marker-types": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/material-amount-units-supported (PWG5100.21)
	"material-amount-units-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/material-diameter-supported (PWG5100.21)
	"material-diameter-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/material-nozzle-diameter-supported (PWG5100.21)
	"material-nozzle-diameter-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/material-purpose-supported (PWG5100.21)
	"material-purpose-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/material-rate-supported (PWG5100.21)
	"material-rate-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/material-rate-units-supported (PWG5100.21)
	"material-rate-units-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/material-shell-thickness-supported (PWG5100.21)
	"material-shell-thickness-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/material-temperature-supported (PWG5100.21)
	"material-temperature-supported": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/material-type-supported (PWG5100.21)
	"material-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/materials-col-database (PWG5100.21)
	"materials-col-database": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/materials-col-default (PWG5100.21)
	"materials-col-default": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/materials-col-ready (PWG5100.21)
	"materials-col-ready": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/materials-col-supported (PWG5100.21)
	"materials-col-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/max-client-info-supported (PWG5100.7)
	"max-client-info-supported": &DefAttr{
		SetOf: false,
		Min:   4,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/max-materials-col-supported (PWG5100.21)
	"max-materials-col-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/max-page-ranges-supported (PWG5100.7)
	"max-page-ranges-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/max-save-info-supported (PWG5100.11)
	"max-save-info-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/max-stitching-locations-supported (PWG5100.1)
	"max-stitching-locations-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/media-back-coating-supported (PWG5100.7)
	"media-back-coating-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-bottom-margin-supported (PWG5100.7)
	"media-bottom-margin-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/media-col-database (PWG5100.7)
	"media-col-database": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/media-col-database/media-size (PWG5100.7)
			"media-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Printer Description/media-col-database/media-size/x-dimension (PWG5100.7)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
					},
					// Printer Description/media-col-database/media-size/y-dimension (PWG5100.7)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
					},
				}},
			},
			// Printer Description/media-col-database/media-source-properties (PWG5100.7)
			"media-source-properties": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Printer Description/media-col-database/media-source-properties/media-source-feed-direction (PWG5100.7)
					"media-source-feed-direction": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Printer Description/media-col-database/media-source-properties/media-source-feed-orientation (PWG5100.7)
					"media-source-feed-orientation": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagEnum},
					},
				}},
			},
		}},
	},
	// Printer Description/media-col-default (PWG5100.7)
	"media-col-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/media-col-ready (PWG5100.7)
	"media-col-ready": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/media-col-ready/media-size (PWG5100.7)
			"media-size": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Printer Description/media-col-ready/media-size/x-dimension (PWG5100.7)
					"x-dimension": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
					// Printer Description/media-col-ready/media-size/y-dimension (PWG5100.7)
					"y-dimension": &DefAttr{
						SetOf: false,
						Min:   0,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagInteger},
					},
				}},
			},
			// Printer Description/media-col-ready/media-source-properties (PWG5100.7)
			"media-source-properties": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
				Members: []map[string]*DefAttr{{
					// Printer Description/media-col-ready/media-source-properties/media-source-feed-direction (PWG5100.7)
					"media-source-feed-direction": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   255,
						Tags:  []goipp.Tag{goipp.TagKeyword},
					},
					// Printer Description/media-col-ready/media-source-properties/media-source-feed-orientation (PWG5100.7)
					"media-source-feed-orientation": &DefAttr{
						SetOf: false,
						Min:   1,
						Max:   MAX,
						Tags:  []goipp.Tag{goipp.TagEnum},
					},
				}},
			},
		}},
	},
	// Printer Description/media-col-supported (PWG5100.7)
	"media-col-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/media-color-supported (PWG5100.7)
	"media-color-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-default (rfc8011)
	"media-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/media-front-coating-supported (PWG5100.7)
	"media-front-coating-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-grain-supported (PWG5100.7)
	"media-grain-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-hole-count-supported (PWG5100.7)
	"media-hole-count-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/media-info-supported (PWG5100.7)
	"media-info-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/media-key-supported (PWG5100.7)
	"media-key-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-left-margin-supported (PWG5100.7)
	"media-left-margin-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/media-order-count-supported (PWG5100.7)
	"media-order-count-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/media-overprint-default (PWG5100.13)
	"media-overprint-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
	},
	// Printer Description/media-overprint-distance-supported (PWG5100.13)
	"media-overprint-distance-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/media-overprint-method-supported (PWG5100.13)
	"media-overprint-method-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/media-overprint-supported (PWG5100.13)
	"media-overprint-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/media-pre-printed-supported (PWG5100.7)
	"media-pre-printed-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-ready (rfc8011)
	"media-ready": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-recycled-supported (PWG5100.7)
	"media-recycled-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-right-margin-supported (PWG5100.7)
	"media-right-margin-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/media-size-supported (PWG5100.7)
	"media-size-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/media-size-supported/x-dimension (PWG5100.7)
			"x-dimension": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
			},
			// Printer Description/media-size-supported/y-dimension (PWG5100.7)
			"y-dimension": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
			},
		}},
	},
	// Printer Description/media-source-supported (PWG5100.7)
	"media-source-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-supported (rfc8011)
	"media-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-thickness-supported (PWG5100.7)
	"media-thickness-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/media-tooth-supported (PWG5100.7)
	"media-tooth-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-top-margin-supported (PWG5100.7)
	"media-top-margin-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/media-top-offset-supported (IPPLABEL)
	"media-top-offset-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   -2147483648,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/media-tracking-supported (IPPLABEL)
	"media-tracking-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/media-type-supported (PWG5100.7)
	"media-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/media-weight-metric-supported (PWG5100.7)
	"media-weight-metric-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/message-supported (PWG5100.15)
	"message-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/multiple-destination-uris-supported (PWG5100.15)
	"multiple-destination-uris-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/multiple-document-handling-default (rfc8011)
	"multiple-document-handling-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/multiple-document-handling-supported (rfc8011)
	"multiple-document-handling-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/multiple-document-jobs-supported (rfc8011)
	"multiple-document-jobs-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/multiple-object-handling-default (PWG5100.21)
	"multiple-object-handling-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/multiple-object-handling-supported (PWG5100.21)
	"multiple-object-handling-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/multiple-operation-time-out (rfc8011)
	"multiple-operation-time-out": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/multiple-operation-time-out-action (PWG5100.13)
	"multiple-operation-time-out-action": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/natural-language-configured (rfc8011)
	"natural-language-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Printer Description/notify-attributes-supported (rfc3995)
	"notify-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/notify-events-default (rfc3995)
	"notify-events-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/notify-events-supported (rfc3995)
	"notify-events-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/notify-lease-duration-default (rfc3995)
	"notify-lease-duration-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   67108863,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/notify-lease-duration-supported (rfc3995)
	"notify-lease-duration-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   67108863,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/notify-pull-method-supported (rfc3995)
	"notify-pull-method-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/notify-schemes-supported (rfc3995)
	"notify-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/number-of-retries-default (PWG5100.15)
	"number-of-retries-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/number-of-retries-supported (PWG5100.15)
	"number-of-retries-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/number-up-default (rfc8011)
	"number-up-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/number-up-supported (rfc8011)
	"number-up-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/oauth-authorization-scope (PWG5100.23)
	"oauth-authorization-scope": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/oauth-authorization-server-uri (PWG5100.23)
	"oauth-authorization-server-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// Printer Description/operations-supported (rfc8011)
	"operations-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/organization-name-supported (PWG5100.15)
	"organization-name-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/orientation-requested-default (rfc8011)
	"orientation-requested-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum, goipp.TagNoValue},
	},
	// Printer Description/orientation-requested-supported (rfc8011)
	"orientation-requested-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/output-attributes-default (PWG5100.17)
	"output-attributes-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
	},
	// Printer Description/output-attributes-supported (PWG5100.17)
	"output-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/output-bin-default (PWG5100.2)
	"output-bin-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/output-bin-supported (PWG5100.2)
	"output-bin-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/output-device-supported (PWG5100.7)
	"output-device-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/output-device-uuid-supported (PWG5100.18)
	"output-device-uuid-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/overrides-supported (PWG5100.6)
	"overrides-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/page-delivery-default (PWG5100.3)
	"page-delivery-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/page-delivery-supported (PWG5100.3)
	"page-delivery-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/page-order-received-default (PWG5100.3)
	"page-order-received-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/page-order-received-supported (PWG5100.3)
	"page-order-received-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/page-ranges-supported (rfc8011)
	"page-ranges-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/pages-per-subset-supported (PWG5100.13)
	"pages-per-subset-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/parent-printers-supported (rfc3998)
	"parent-printers-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/pclm-raster-back-side (HP20180907)
	"pclm-raster-back-side": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pclm-source-resolution-default (HP20180907)
	"pclm-source-resolution-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/pclm-source-resolution-supported (HP20180907)
	"pclm-source-resolution-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/pclm-strip-height-preferred (HP20180907)
	"pclm-strip-height-preferred": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/pclm-strip-height-supported (HP20180907)
	"pclm-strip-height-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/pdf-features-supported (PWG5100.21)
	"pdf-features-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pdf-k-octets-supported (PWG5100.13)
	"pdf-k-octets-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/pdf-versions-supported (PWG5100.13)
	"pdf-versions-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pdl-init-file-entry-supported (PWG5100.11)
	"pdl-init-file-entry-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/pdl-init-file-location-supported (PWG5100.11)
	"pdl-init-file-location-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/pdl-init-file-name-subdirectory-supported (PWG5100.11)
	"pdl-init-file-name-subdirectory-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/pdl-init-file-name-supported (PWG5100.11)
	"pdl-init-file-name-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/pdl-init-file-supported (PWG5100.11)
	"pdl-init-file-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pdl-override-guaranteed-supported (IPPWG20151019)
	"pdl-override-guaranteed-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pdl-override-supported (rfc8011)
	"pdl-override-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pkcs7-document-format-supported (PWG5100.TRUSTNOONE)
	"pkcs7-document-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/platform-shape (PWG5100.21)
	"platform-shape": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/platform-temperature-default (PWG5100.21)
	"platform-temperature-default": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/platform-temperature-supported (PWG5100.21)
	"platform-temperature-supported": &DefAttr{
		SetOf: true,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/port-monitor (CUPS)
	"port-monitor": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/port-monitor-supported (CUPS)
	"port-monitor-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/ppd-name (CUPS)
	"ppd-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/preferred-attributes-supported (PWG5100.13)
	"preferred-attributes-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/presentation-direction-number-up-default (PWG5100.3)
	"presentation-direction-number-up-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/presentation-direction-number-up-supported (PWG5100.3)
	"presentation-direction-number-up-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-accuracy-supported (PWG5100.21)
	"print-accuracy-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/print-accuracy-supported/accuracy-units (PWG5100.21)
			"accuracy-units": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/print-accuracy-supported/x-accuracy (PWG5100.21)
			"x-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Printer Description/print-accuracy-supported/y-accuracy (PWG5100.21)
			"y-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Printer Description/print-accuracy-supported/z-accuracy (PWG5100.21)
			"z-accuracy": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Printer Description/print-base-default (PWG5100.21)
	"print-base-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-base-supported (PWG5100.21)
	"print-base-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-color-mode-default (PWG5100.13)
	"print-color-mode-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-color-mode-icc-profiles (PWG5100.13)
	"print-color-mode-icc-profiles": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/print-color-mode-icc-profiles/print-color-mode (PWG5100.13)
			"print-color-mode": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/print-color-mode-icc-profiles/profile-uri (PWG5100.13)
			"profile-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Printer Description/print-color-mode-supported (PWG5100.13)
	"print-color-mode-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-content-optimize-default (PWG5100.7)
	"print-content-optimize-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-content-optimize-supported (PWG5100.7)
	"print-content-optimize-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-darkness-default (IPPLABEL)
	"print-darkness-default": &DefAttr{
		SetOf: false,
		Min:   -100,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/print-darkness-supported (IPPLABEL)
	"print-darkness-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/print-objects-supported (PWG5100.21)
	"print-objects-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-processing-attributes-supported (PWG5100.13)
	"print-processing-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-quality-default (rfc8011)
	"print-quality-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/print-quality-supported (rfc8011)
	"print-quality-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/print-rendering-intent-default (PWG5100.13)
	"print-rendering-intent-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-rendering-intent-supported (PWG5100.13)
	"print-rendering-intent-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-scaling-default (PWG5100.13)
	"print-scaling-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-scaling-supported (PWG5100.13)
	"print-scaling-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-speed-default (IPPLABEL)
	"print-speed-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/print-speed-supported (IPPLABEL)
	"print-speed-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/print-supports-default (PWG5100.21)
	"print-supports-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/print-supports-supported (PWG5100.21)
	"print-supports-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-asset-tag (PWG5100.11)
	"printer-asset-tag": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Description/printer-camera-image-uri (PWG5100.21)
	"printer-camera-image-uri": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-charge-info (PWG5100.16)
	"printer-charge-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-charge-info-uri (PWG5100.16)
	"printer-charge-info-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-commands (CUPS)
	"printer-commands": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-contact-col (PWG5100.22)
	"printer-contact-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagUnknown},
		Members: []map[string]*DefAttr{{
			// Printer Description/printer-contact-col/contact-name (PWG5100.22)
			"contact-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Printer Description/printer-contact-col/contact-uri (PWG5100.22)
			"contact-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Printer Description/printer-contact-col/contact-vcard (PWG5100.22)
			"contact-vcard": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
		}},
	},
	// Printer Description/printer-current-time (rfc8011)
	"printer-current-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagUnknown},
	},
	// Printer Description/printer-darkness-configured (IPPLABEL)
	"printer-darkness-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/printer-darkness-supported (IPPLABEL)
	"printer-darkness-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/printer-detailed-status-messages (PWG5100.11)
	"printer-detailed-status-messages": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-device-id (PWG5107.2)
	"printer-device-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-dns-sd-name (PWG5100.13)
	"printer-dns-sd-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// Printer Description/printer-driver-installer (rfc8011)
	"printer-driver-installer": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-error-policy (CUPS)
	"printer-error-policy": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-error-policy-supported (CUPS)
	"printer-error-policy-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-fax-log-uri (PWG5100.15)
	"printer-fax-log-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-fax-modem-info (PWG5100.15)
	"printer-fax-modem-info": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-fax-modem-name (PWG5100.15)
	"printer-fax-modem-name": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/printer-fax-modem-number (PWG5100.15)
	"printer-fax-modem-number": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-geo-location (PWG5100.13)
	"printer-geo-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagUnknown},
	},
	// Printer Description/printer-get-attributes-supported (PWG5100.13)
	"printer-get-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-icc-profiles (PWG5100.13)
	"printer-icc-profiles": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/printer-icc-profiles/profile-name (PWG5100.13)
			"profile-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Printer Description/printer-icc-profiles/profile-url (PWG5100.13)
			"profile-url": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Printer Description/printer-icons (PWG5100.13)
	"printer-icons": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-id (CUPS)
	"printer-id": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/printer-info (rfc8011)
	"printer-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-is-accepting-jobs (CUPS)
	"printer-is-accepting-jobs": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/printer-is-shared (CUPS)
	"printer-is-shared": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/printer-is-temporary (CUPS)
	"printer-is-temporary": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/printer-kind (PWG5100.13)
	"printer-kind": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/printer-location (rfc8011)
	"printer-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-make-and-model (rfc8011)
	"printer-make-and-model": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-mandatory-job-attributes (PWG5100.13)
	"printer-mandatory-job-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-mode-configured (PWG5100.18)
	"printer-mode-configured": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-mode-supported (PWG5100.18)
	"printer-mode-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-more-info (CUPS)
	"printer-more-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-more-info-manufacturer (rfc8011)
	"printer-more-info-manufacturer": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-name (rfc8011)
	"printer-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/printer-op-policy (CUPS)
	"printer-op-policy": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/printer-op-policy-supported (CUPS)
	"printer-op-policy-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/printer-organization (PWG5100.13)
	"printer-organization": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-organizational-unit (PWG5100.13)
	"printer-organizational-unit": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-pkcs7-public-key (PWG5100.TRUSTNOONE)
	"printer-pkcs7-public-key": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-pkcs7-repertoire-configured (PWG5100.TRUSTNOONE)
	"printer-pkcs7-repertoire-configured": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-pkcs7-repertoire-supported (PWG5100.TRUSTNOONE)
	"printer-pkcs7-repertoire-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-privacy-policy-uri (IPPPRIVACY10)
	"printer-privacy-policy-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-requested-job-attributes (PWG5100.16)
	"printer-requested-job-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/printer-resolution-default (rfc8011)
	"printer-resolution-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/printer-resolution-supported (rfc8011)
	"printer-resolution-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/printer-service-contact-col (PWG5100.11)
	"printer-service-contact-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/printer-service-contact-col/contact-name (PWG5100.11)
			"contact-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// Printer Description/printer-service-contact-col/contact-uri (PWG5100.11)
			"contact-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
			// Printer Description/printer-service-contact-col/contact-vcard (PWG5100.11)
			"contact-vcard": &DefAttr{
				SetOf: true,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagText},
			},
		}},
	},
	// Printer Description/printer-state (CUPS)
	"printer-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/printer-state-message (CUPS)
	"printer-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/printer-static-resource-directory-uri (PWG5100.18)
	"printer-static-resource-directory-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/printer-static-resource-k-octets-supported (PWG5100.18)
	"printer-static-resource-k-octets-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/printer-strings-languages-supported (PWG5100.13)
	"printer-strings-languages-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Printer Description/printer-strings-uri (PWG5100.13)
	"printer-strings-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// Printer Description/printer-type (CUPS)
	"printer-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/printer-type-mask (CUPS)
	"printer-type-mask": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Description/printer-volume-supported (PWG5100.21)
	"printer-volume-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/printer-volume-supported/x-dimension (PWG5100.21)
			"x-dimension": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Printer Description/printer-volume-supported/y-dimension (PWG5100.21)
			"y-dimension": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// Printer Description/printer-volume-supported/z-dimension (PWG5100.21)
			"z-dimension": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// Printer Description/printer-wifi-password (IPPWIFI)
	"printer-wifi-password": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Description/printer-wifi-ssid (IPPWIFI)
	"printer-wifi-ssid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/printer-xri-supported (rfc3380)
	"printer-xri-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// Printer Description/printer-xri-supported/xri-authentication (rfc3380)
			"xri-authentication": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/printer-xri-supported/xri-security (rfc3380)
			"xri-security": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// Printer Description/printer-xri-supported/xri-uri (rfc3380)
			"xri-uri": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   1023,
				Tags:  []goipp.Tag{goipp.TagURI},
			},
		}},
	},
	// Printer Description/proof-copies-supported (PWG5100.11)
	"proof-copies-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/proof-print-copies-supported (PWG5100.11)
	"proof-print-copies-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/proof-print-default (PWG5100.11)
	"proof-print-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagNoValue},
	},
	// Printer Description/proof-print-supported (PWG5100.11)
	"proof-print-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/punching-hole-diameter-configured (PWG5100.1)
	"punching-hole-diameter-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/punching-locations-supported (PWG5100.1)
	"punching-locations-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/punching-offset-supported (PWG5100.1)
	"punching-offset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/punching-reference-edge-supported (PWG5100.1)
	"punching-reference-edge-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pwg-raster-document-resolution-supported (PWG5102.4)
	"pwg-raster-document-resolution-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagResolution},
	},
	// Printer Description/pwg-raster-document-sheet-back (PWG5102.4)
	"pwg-raster-document-sheet-back": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pwg-raster-document-type-supported (PWG5102.4)
	"pwg-raster-document-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/pwg-safe-gcode-supported (PWG5199.7)
	"pwg-safe-gcode-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Description/reference-uri-schemes-supported (rfc8011)
	"reference-uri-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/repertoire-supported (PWG5101.2)
	"repertoire-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// Printer Description/requesting-user-name-allowed (CUPS)
	"requesting-user-name-allowed": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagDeleteAttr},
	},
	// Printer Description/requesting-user-name-denied (CUPS)
	"requesting-user-name-denied": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagDeleteAttr},
	},
	// Printer Description/requesting-user-uri-schemes-supported (PWG5100.13)
	"requesting-user-uri-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// Printer Description/requesting-user-uri-supported (PWG5100.13)
	"requesting-user-uri-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/retry-interval-default (PWG5100.15)
	"retry-interval-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/retry-interval-supported (PWG5100.15)
	"retry-interval-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/retry-time-out-default (PWG5100.15)
	"retry-time-out-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/retry-time-out-supported (PWG5100.15)
	"retry-time-out-supported": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/save-disposition-supported (PWG5100.11)
	"save-disposition-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/save-document-format-default (PWG5100.11)
	"save-document-format-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/save-document-format-supported (PWG5100.11)
	"save-document-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/save-location-default (PWG5100.11)
	"save-location-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/save-location-supported (PWG5100.11)
	"save-location-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/save-name-subdirectory-supported (PWG5100.11)
	"save-name-subdirectory-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/save-name-supported (PWG5100.11)
	"save-name-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Description/separator-sheets-default (PWG5100.3)
	"separator-sheets-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Description/separator-sheets-supported (PWG5100.3)
	"separator-sheets-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/sheet-collate-default (rfc3381)
	"sheet-collate-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/sheet-collate-supported (rfc3381)
	"sheet-collate-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/sides-default (rfc8011)
	"sides-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/sides-supported (rfc8011)
	"sides-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/smi2699-auth-print-group (IPPSERVER)
	"smi2699-auth-print-group": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/smi2699-auth-proxy-group (IPPSERVER)
	"smi2699-auth-proxy-group": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/smi2699-device-command (IPPSERVER)
	"smi2699-device-command": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/smi2699-device-format (IPPSERVER)
	"smi2699-device-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Printer Description/smi2699-device-name (IPPSERVER)
	"smi2699-device-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Description/smi2699-device-uri (IPPSERVER)
	"smi2699-device-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/stitching-angle-supported (PWG5100.1)
	"stitching-angle-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   359,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/stitching-locations-supported (PWG5100.1)
	"stitching-locations-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/stitching-method-supported (PWG5100.1)
	"stitching-method-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/stitching-offset-supported (PWG5100.1)
	"stitching-offset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/stitching-reference-edge-supported (PWG5100.1)
	"stitching-reference-edge-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/subject-supported (PWG5100.15)
	"subject-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/subordinate-printers-supported (rfc3998)
	"subordinate-printers-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Description/subscription-privacy-attributes (IPPPRIVACY10)
	"subscription-privacy-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/subscription-privacy-scope (IPPPRIVACY10)
	"subscription-privacy-scope": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/to-name-supported (PWG5100.15)
	"to-name-supported": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/trimming-offset-supported (PWG5100.1)
	"trimming-offset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// Printer Description/trimming-reference-edge-supported (PWG5100.1)
	"trimming-reference-edge-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/trimming-type-supported (PWG5100.1)
	"trimming-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/trimming-when-supported (PWG5100.1)
	"trimming-when-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/urf-supported (CUPS)
	"urf-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/uri-authentication-supported (rfc8011)
	"uri-authentication-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/uri-security-supported (rfc8011)
	"uri-security-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/user-defined-values-supported (PWG5100.3)
	"user-defined-values-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/which-jobs-supported (PWG5100.7)
	"which-jobs-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/x-image-position-default (PWG5100.3)
	"x-image-position-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/x-image-position-supported (PWG5100.3)
	"x-image-position-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/x-image-shift-default (PWG5100.3)
	"x-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/x-image-shift-supported (PWG5100.3)
	"x-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/x-side1-image-shift-default (PWG5100.3)
	"x-side1-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/x-side1-image-shift-supported (PWG5100.3)
	"x-side1-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/x-side2-image-shift-default (PWG5100.3)
	"x-side2-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/x-side2-image-shift-supported (PWG5100.3)
	"x-side2-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/y-image-position-default (PWG5100.3)
	"y-image-position-default": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/y-image-position-supported (PWG5100.3)
	"y-image-position-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Description/y-image-shift-default (PWG5100.3)
	"y-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/y-image-shift-supported (PWG5100.3)
	"y-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/y-side1-image-shift-default (PWG5100.3)
	"y-side1-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/y-side1-image-shift-supported (PWG5100.3)
	"y-side1-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
	// Printer Description/y-side2-image-shift-default (PWG5100.3)
	"y-side2-image-shift-default": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Description/y-side2-image-shift-supported (PWG5100.3)
	"y-side2-image-shift-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagRange},
	},
}

// PrinterStatus is the Printer Status attributes
var PrinterStatus = map[string]*DefAttr{
	// Printer Status/chamber-humidity-current (PWG5100.21)
	"chamber-humidity-current": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   100,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagUnknown},
	},
	// Printer Status/chamber-temperature-current (PWG5100.21)
	"chamber-temperature-current": &DefAttr{
		SetOf: false,
		Min:   -273,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagUnknown},
	},
	// Printer Status/device-service-count (PWG5100.13)
	"device-service-count": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/device-uuid (PWG5100.13)
	"device-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/document-format-varying-attributes (rfc3380)
	"document-format-varying-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/job-settable-attributes-supported (rfc3380)
	"job-settable-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/pages-per-minute (rfc8011)
	"pages-per-minute": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/pages-per-minute-color (rfc8011)
	"pages-per-minute-color": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-alert (PWG5100.9)
	"printer-alert": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-alert-description (PWG5100.9)
	"printer-alert-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-camera-image-uri (PWG5100.21)
	"printer-camera-image-uri": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-config-change-date-time (PWG5100.13)
	"printer-config-change-date-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagUnknown},
	},
	// Printer Status/printer-config-change-time (PWG5100.13)
	"printer-config-change-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-config-changes (PWG5100.22)
	"printer-config-changes": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-console-display (IPPCONSOLE)
	"printer-console-display": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-console-light (IPPCONSOLE)
	"printer-console-light": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-console-light-description (IPPCONSOLE)
	"printer-console-light-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-cover (IPP20210223)
	"printer-cover": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-cover-description (IPP20210223)
	"printer-cover-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-detailed-status-messages (PWG5100.7)
	"printer-detailed-status-messages": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-finisher (PWG5100.1)
	"printer-finisher": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-finisher-description (PWG5100.1)
	"printer-finisher-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-finisher-supplies (PWG5100.1)
	"printer-finisher-supplies": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-finisher-supplies-description (PWG5100.1)
	"printer-finisher-supplies-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-firmware-name (PWG5100.13)
	"printer-firmware-name": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Printer Status/printer-firmware-patches (PWG5100.13)
	"printer-firmware-patches": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-firmware-string-version (PWG5100.13)
	"printer-firmware-string-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-firmware-version (PWG5100.13)
	"printer-firmware-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-home-page-uri (IPPCONSOLE)
	"printer-home-page-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-id (PWG5100.22)
	"printer-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-impressions-completed (PWG5100.22)
	"printer-impressions-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-impressions-completed-col (PWG5100.22)
	"printer-impressions-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Status/printer-input-tray (PWG5100.13)
	"printer-input-tray": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-is-accepting-jobs (rfc8011)
	"printer-is-accepting-jobs": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// Printer Status/printer-media-sheets-completed (PWG5100.22)
	"printer-media-sheets-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-media-sheets-completed-col (PWG5100.22)
	"printer-media-sheets-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Status/printer-message-date-time (rfc3380)
	"printer-message-date-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Printer Status/printer-message-from-operator (rfc8011)
	"printer-message-from-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-message-time (rfc3380)
	"printer-message-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-more-info (rfc8011)
	"printer-more-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-output-tray (PWG5100.13)
	"printer-output-tray": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-pages-completed (PWG5100.22)
	"printer-pages-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-pages-completed-col (PWG5100.22)
	"printer-pages-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// Printer Status/printer-serial-number (PWG5100.11)
	"printer-serial-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-service-type (PWG5100.22)
	"printer-service-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/printer-settable-attributes-supported (rfc3380)
	"printer-settable-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/printer-state (rfc8011)
	"printer-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Status/printer-state-change-date-time (rfc3995)
	"printer-state-change-date-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Printer Status/printer-state-change-time (rfc3995)
	"printer-state-change-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-state-message (rfc8011)
	"printer-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-state-reasons (rfc8011)
	"printer-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/printer-static-resource-k-octets-free (PWG5100.18)
	"printer-static-resource-k-octets-free": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-storage (PWG5100.11)
	"printer-storage": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-storage-description (PWG5100.11)
	"printer-storage-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-supply (PWG5100.13)
	"printer-supply": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Printer Status/printer-supply-description (PWG5100.13)
	"printer-supply-description": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Printer Status/printer-supply-info-uri (PWG5100.13)
	"printer-supply-info-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-up-time (rfc8011)
	"printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/printer-uri-supported (rfc8011)
	"printer-uri-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-uuid (PWG5100.13)
	"printer-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Printer Status/printer-wifi-state (IPPWIFI)
	"printer-wifi-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Printer Status/queued-job-count (rfc8011)
	"queued-job-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Printer Status/xri-authentication-supported (rfc3380)
	"xri-authentication-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/xri-security-supported (rfc3380)
	"xri-security-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Printer Status/xri-uri-scheme-supported (rfc3380)
	"xri-uri-scheme-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
}

// ResourceDescription is the Resource Description attributes
var ResourceDescription = map[string]*DefAttr{
	// Resource Description/resource-info (PWG5100.22)
	"resource-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Resource Description/resource-name (PWG5100.22)
	"resource-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
}

// ResourceStatus is the Resource Status attributes
var ResourceStatus = map[string]*DefAttr{
	// Resource Status/date-time-at-canceled (PWG5100.22)
	"date-time-at-canceled": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagNoValue},
	},
	// Resource Status/date-time-at-creation (PWG5100.22)
	"date-time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// Resource Status/date-time-at-installed (PWG5100.22)
	"date-time-at-installed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime, goipp.TagNoValue},
	},
	// Resource Status/resource-data-uri (PWG5100.22)
	"resource-data-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// Resource Status/resource-format (PWG5100.22)
	"resource-format": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// Resource Status/resource-id (PWG5100.22)
	"resource-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Resource Status/resource-k-octets (PWG5100.22)
	"resource-k-octets": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Resource Status/resource-natural-language (PWG5100.22)
	"resource-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Resource Status/resource-patches (PWG5100.22)
	"resource-patches": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText, goipp.TagNoValue},
	},
	// Resource Status/resource-signature (PWG5100.22)
	"resource-signature": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// Resource Status/resource-state (PWG5100.22)
	"resource-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Resource Status/resource-state-message (PWG5100.22)
	"resource-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// Resource Status/resource-state-reasons (PWG5100.22)
	"resource-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Resource Status/resource-string-version (PWG5100.22)
	"resource-string-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText, goipp.TagNoValue},
	},
	// Resource Status/resource-type (PWG5100.22)
	"resource-type": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Resource Status/resource-use-count (PWG5100.22)
	"resource-use-count": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Resource Status/resource-uuid (PWG5100.22)
	"resource-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   45,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Resource Status/resource-version (PWG5100.22)
	"resource-version": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   64,
		Tags:  []goipp.Tag{goipp.TagString, goipp.TagNoValue},
	},
	// Resource Status/time-at-canceled (PWG5100.22)
	"time-at-canceled": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// Resource Status/time-at-creation (PWG5100.22)
	"time-at-creation": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Resource Status/time-at-installed (PWG5100.22)
	"time-at-installed": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
}

// SubscriptionStatus is the Subscription Status attributes
var SubscriptionStatus = map[string]*DefAttr{
	// Subscription Status/notify-job-id (rfc3995)
	"notify-job-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-lease-expiration-time (rfc3995)
	"notify-lease-expiration-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-printer-up-time (rfc3995)
	"notify-printer-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-printer-uri (rfc3995)
	"notify-printer-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Subscription Status/notify-resource-id (PWG5100.22)
	"notify-resource-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-sequence-number (rfc3995)
	"notify-sequence-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-status-code (rfc3995)
	"notify-status-code": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// Subscription Status/notify-subscriber-user-name (rfc3995)
	"notify-subscriber-user-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// Subscription Status/notify-subscriber-user-uri (PWG5100.13)
	"notify-subscriber-user-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Subscription Status/notify-subscription-id (rfc3995)
	"notify-subscription-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-subscription-uuid (PWG5100.13)
	"notify-subscription-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Subscription Status/notify-system-up-time (PWG5100.22)
	"notify-system-up-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Status/notify-system-uri (PWG5100.22)
	"notify-system-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
}

// SubscriptionTemplate is the Subscription Template attributes
var SubscriptionTemplate = map[string]*DefAttr{
	// Subscription Template/notify-attributes (rfc3995)
	"notify-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Subscription Template/notify-charset (rfc3995)
	"notify-charset": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// Subscription Template/notify-events (rfc3995)
	"notify-events": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Subscription Template/notify-lease-duration (rfc3995)
	"notify-lease-duration": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   67108863,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Template/notify-max-events-supported (rfc3995)
	"notify-max-events-supported": &DefAttr{
		SetOf: false,
		Min:   2,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Template/notify-natural-language (rfc3995)
	"notify-natural-language": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// Subscription Template/notify-pull-method (rfc3995)
	"notify-pull-method": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// Subscription Template/notify-recipient-uri (rfc3995)
	"notify-recipient-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// Subscription Template/notify-time-interval (rfc3995)
	"notify-time-interval": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// Subscription Template/notify-user-data (rfc3995)
	"notify-user-data": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagString},
	},
}

// SystemDescription is the System Description attributes
var SystemDescription = map[string]*DefAttr{
	// System Description/charset-configured (PWG5100.22)
	"charset-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// System Description/charset-supported (PWG5100.22)
	"charset-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagCharset},
	},
	// System Description/document-format-supported (PWG5100.22)
	"document-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// System Description/generated-natural-language-supported (PWG5100.22)
	"generated-natural-language-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// System Description/ipp-features-supported (PWG5100.22)
	"ipp-features-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/ipp-versions-supported (PWG5100.22)
	"ipp-versions-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/ippget-event-life (PWG5100.22)
	"ippget-event-life": &DefAttr{
		SetOf: false,
		Min:   15,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Description/multiple-document-printers-supported (PWG5100.22)
	"multiple-document-printers-supported": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBoolean},
	},
	// System Description/natural-language-configured (PWG5100.22)
	"natural-language-configured": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// System Description/notify-attributes-supported (PWG5100.22)
	"notify-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/notify-events-default (PWG5100.22)
	"notify-events-default": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/notify-events-supported (PWG5100.22)
	"notify-events-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/notify-lease-duration-default (PWG5100.22)
	"notify-lease-duration-default": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   67108863,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Description/notify-lease-duration-supported (PWG5100.22)
	"notify-lease-duration-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   67108863,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagRange},
	},
	// System Description/notify-max-events-supported (PWG5100.22)
	"notify-max-events-supported": &DefAttr{
		SetOf: false,
		Min:   2,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Description/notify-pull-method-supported (PWG5100.22)
	"notify-pull-method-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/notify-schemes-supported (PWG5100.22)
	"notify-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// System Description/oauth-authorization-scope (PWG5100.23)
	"oauth-authorization-scope": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName, goipp.TagNoValue},
	},
	// System Description/oauth-authorization-server-uri (PWG5100.23)
	"oauth-authorization-server-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// System Description/operations-supported (PWG5100.22)
	"operations-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// System Description/output-device-x509-type-supported (PWG5100.22)
	"output-device-x509-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/power-calendar-policy-col (PWG5100.22)
	"power-calendar-policy-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Description/power-calendar-policy-col/calendar-id (PWG5100.22)
			"calendar-id": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/day-of-month (PWG5100.22)
			"day-of-month": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   31,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/day-of-week (PWG5100.22)
			"day-of-week": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   7,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/hour (PWG5100.22)
			"hour": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   23,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/minute (PWG5100.22)
			"minute": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   59,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/month (PWG5100.22)
			"month": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   12,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-calendar-policy-col/request-power-state (PWG5100.22)
			"request-power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Description/power-calendar-policy-col/run-once (PWG5100.22)
			"run-once": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
		}},
	},
	// System Description/power-event-policy-col (PWG5100.22)
	"power-event-policy-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Description/power-event-policy-col/event-id (PWG5100.22)
			"event-id": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-event-policy-col/event-name (PWG5100.22)
			"event-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// System Description/power-event-policy-col/request-power-state (PWG5100.22)
			"request-power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// System Description/power-timeout-policy-col (PWG5100.22)
	"power-timeout-policy-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Description/power-timeout-policy-col/start-power-state (PWG5100.22)
			"start-power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Description/power-timeout-policy-col/timeout-id (PWG5100.22)
			"timeout-id": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Description/power-timeout-policy-col/timeout-predicate (PWG5100.22)
			"timeout-predicate": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Description/power-timeout-policy-col/timeout-seconds (PWG5100.22)
			"timeout-seconds": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// System Description/printer-creation-attributes-supported (PWG5100.22)
	"printer-creation-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/resource-format-supported (PWG5100.22)
	"resource-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// System Description/resource-settable-attributes-supported (PWG5100.22)
	"resource-settable-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/resource-type-supported (PWG5100.22)
	"resource-type-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/smi2699-auth-group-supported (IPPSERVER)
	"smi2699-auth-group-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Description/smi2699-device-command-supported (IPPSERVER)
	"smi2699-device-command-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Description/smi2699-device-format-supported (IPPSERVER)
	"smi2699-device-format-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagMimeType},
	},
	// System Description/smi2699-device-uri-schemes-supported (IPPSERVER)
	"smi2699-device-uri-schemes-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
	// System Description/system-asset-tag (PWG5100.22)
	"system-asset-tag": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// System Description/system-contact-col (PWG5100.22)
	"system-contact-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagUnknown},
	},
	// System Description/system-current-time (PWG5100.22)
	"system-current-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// System Description/system-default-printer-id (PWG5100.22)
	"system-default-printer-id": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   65535,
		Tags:  []goipp.Tag{goipp.TagInteger, goipp.TagNoValue},
	},
	// System Description/system-dns-sd-name (PWG5100.22)
	"system-dns-sd-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Description/system-geo-location (PWG5100.22)
	"system-geo-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagUnknown},
	},
	// System Description/system-info (PWG5100.22)
	"system-info": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Description/system-location (PWG5100.22)
	"system-location": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Description/system-make-and-model (PWG5100.22)
	"system-make-and-model": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Description/system-mandatory-printer-attributes (PWG5100.22)
	"system-mandatory-printer-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/system-mandatory-registration-attributes (PWG5100.22)
	"system-mandatory-registration-attributes": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/system-message-from-operator (PWG5100.22)
	"system-message-from-operator": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Description/system-name (PWG5100.22)
	"system-name": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   127,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Description/system-service-contact-col (PWG5100.22)
	"system-service-contact-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection, goipp.TagUnknown},
	},
	// System Description/system-settable-attributes-supported (PWG5100.22)
	"system-settable-attributes-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Description/system-strings-languages-supported (PWG5100.22)
	"system-strings-languages-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagLanguage},
	},
	// System Description/system-strings-uri (PWG5100.22)
	"system-strings-uri": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagURI, goipp.TagNoValue},
	},
	// System Description/system-xri-supported (PWG5100.22)
	"system-xri-supported": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
}

// SystemStatus is the System Status attributes
var SystemStatus = map[string]*DefAttr{
	// System Status/power-log-col (PWG5100.22)
	"power-log-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/power-log-col/log-id (PWG5100.22)
			"log-id": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-log-col/power-state (PWG5100.22)
			"power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/power-log-col/power-state-date-time (PWG5100.22)
			"power-state-date-time": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagDateTime},
			},
			// System Status/power-log-col/power-state-message (PWG5100.22)
			"power-state-message": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagText},
			},
		}},
	},
	// System Status/power-state-capabilities-col (PWG5100.22)
	"power-state-capabilities-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/power-state-capabilities-col/can-accept-jobs (PWG5100.22)
			"can-accept-jobs": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// System Status/power-state-capabilities-col/can-process-jobs (PWG5100.22)
			"can-process-jobs": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// System Status/power-state-capabilities-col/power-active-watts (PWG5100.22)
			"power-active-watts": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-capabilities-col/power-inactive-watts (PWG5100.22)
			"power-inactive-watts": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-capabilities-col/power-state (PWG5100.22)
			"power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// System Status/power-state-counters-col (PWG5100.22)
	"power-state-counters-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/power-state-counters-col/hibernate-transitions (PWG5100.22)
			"hibernate-transitions": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-counters-col/on-transitions (PWG5100.22)
			"on-transitions": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-counters-col/standby-transitions (PWG5100.22)
			"standby-transitions": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-counters-col/suspend-transitions (PWG5100.22)
			"suspend-transitions": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// System Status/power-state-monitor-col (PWG5100.22)
	"power-state-monitor-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/power-state-monitor-col/current-month-kwh (PWG5100.22)
			"current-month-kwh": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-monitor-col/current-watts (PWG5100.22)
			"current-watts": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-monitor-col/lifetime-kwh (PWG5100.22)
			"lifetime-kwh": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/power-state-monitor-col/meters-are-actual (PWG5100.22)
			"meters-are-actual": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// System Status/power-state-monitor-col/power-state (PWG5100.22)
			"power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/power-state-monitor-col/power-state-message (PWG5100.22)
			"power-state-message": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// System Status/power-state-monitor-col/power-usage-is-rms-watts (PWG5100.22)
			"power-usage-is-rms-watts": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
		}},
	},
	// System Status/power-state-transitions-col (PWG5100.22)
	"power-state-transitions-col": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/power-state-transitions-col/end-power-state (PWG5100.22)
			"end-power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/power-state-transitions-col/start-power-state (PWG5100.22)
			"start-power-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/power-state-transitions-col/state-transition-seconds (PWG5100.22)
			"state-transition-seconds": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
		}},
	},
	// System Status/system-config-change-date-time (PWG5100.22)
	"system-config-change-date-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// System Status/system-config-change-time (PWG5100.22)
	"system-config-change-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-config-changes (PWG5100.22)
	"system-config-changes": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-configured-printers (PWG5100.22)
	"system-configured-printers": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/system-configured-printers/printer-id (PWG5100.22)
			"printer-id": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   65535,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/system-configured-printers/printer-info (PWG5100.22)
			"printer-info": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// System Status/system-configured-printers/printer-is-accepting-jobs (PWG5100.22)
			"printer-is-accepting-jobs": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBoolean},
			},
			// System Status/system-configured-printers/printer-name (PWG5100.22)
			"printer-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// System Status/system-configured-printers/printer-service-type (PWG5100.22)
			"printer-service-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/system-configured-printers/printer-state (PWG5100.22)
			"printer-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
			// System Status/system-configured-printers/printer-state-reasons (PWG5100.22)
			"printer-state-reasons": &DefAttr{
				SetOf: true,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
			// System Status/system-configured-printers/printer-xri-supported (PWG5100.22)
			"printer-xri-supported": &DefAttr{
				SetOf: false,
				Min:   MIN,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagBeginCollection},
			},
		}},
	},
	// System Status/system-configured-resources (PWG5100.22)
	"system-configured-resources": &DefAttr{
		SetOf: true,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
		Members: []map[string]*DefAttr{{
			// System Status/system-configured-resources/resource-format (PWG5100.22)
			"resource-format": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagMimeType},
			},
			// System Status/system-configured-resources/resource-id (PWG5100.22)
			"resource-id": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagInteger},
			},
			// System Status/system-configured-resources/resource-info (PWG5100.22)
			"resource-info": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagText},
			},
			// System Status/system-configured-resources/resource-name (PWG5100.22)
			"resource-name": &DefAttr{
				SetOf: false,
				Min:   0,
				Max:   127,
				Tags:  []goipp.Tag{goipp.TagName},
			},
			// System Status/system-configured-resources/resource-state (PWG5100.22)
			"resource-state": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   MAX,
				Tags:  []goipp.Tag{goipp.TagEnum},
			},
			// System Status/system-configured-resources/resource-type (PWG5100.22)
			"resource-type": &DefAttr{
				SetOf: false,
				Min:   1,
				Max:   255,
				Tags:  []goipp.Tag{goipp.TagKeyword},
			},
		}},
	},
	// System Status/system-firmware-name (PWG5100.22)
	"system-firmware-name": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Status/system-firmware-patches (PWG5100.22)
	"system-firmware-patches": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-firmware-string-version (PWG5100.22)
	"system-firmware-string-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-firmware-version (PWG5100.22)
	"system-firmware-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   64,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// System Status/system-impressions-completed (PWG5100.22)
	"system-impressions-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-impressions-completed-col (PWG5100.22)
	"system-impressions-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// System Status/system-media-sheets-completed (PWG5100.22)
	"system-media-sheets-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-media-sheets-completed-col (PWG5100.22)
	"system-media-sheets-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// System Status/system-pages-completed (PWG5100.22)
	"system-pages-completed": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-pages-completed-col (PWG5100.22)
	"system-pages-completed-col": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagBeginCollection},
	},
	// System Status/system-resident-application-name (PWG5100.22)
	"system-resident-application-name": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Status/system-resident-application-patches (PWG5100.22)
	"system-resident-application-patches": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-resident-application-string-version (PWG5100.22)
	"system-resident-application-string-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-resident-application-version (PWG5100.22)
	"system-resident-application-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   64,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// System Status/system-serial-number (PWG5100.22)
	"system-serial-number": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-state (PWG5100.22)
	"system-state": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagEnum},
	},
	// System Status/system-state-change-date-time (PWG5100.22)
	"system-state-change-date-time": &DefAttr{
		SetOf: false,
		Min:   MIN,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagDateTime},
	},
	// System Status/system-state-change-time (PWG5100.22)
	"system-state-change-time": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-state-message (PWG5100.22)
	"system-state-message": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-state-reasons (PWG5100.22)
	"system-state-reasons": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Status/system-time-source (PWG5100.22)
	"system-time-source": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword, goipp.TagName},
	},
	// System Status/system-up-time (PWG5100.22)
	"system-up-time": &DefAttr{
		SetOf: false,
		Min:   1,
		Max:   MAX,
		Tags:  []goipp.Tag{goipp.TagInteger},
	},
	// System Status/system-user-application-name (PWG5100.22)
	"system-user-application-name": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagName},
	},
	// System Status/system-user-application-patches (PWG5100.22)
	"system-user-application-patches": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-user-application-string-version (PWG5100.22)
	"system-user-application-string-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   1023,
		Tags:  []goipp.Tag{goipp.TagText},
	},
	// System Status/system-user-application-version (PWG5100.22)
	"system-user-application-version": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   64,
		Tags:  []goipp.Tag{goipp.TagString},
	},
	// System Status/system-uuid (PWG5100.22)
	"system-uuid": &DefAttr{
		SetOf: false,
		Min:   0,
		Max:   45,
		Tags:  []goipp.Tag{goipp.TagURI},
	},
	// System Status/xri-authentication-supported (PWG5100.22)
	"xri-authentication-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Status/xri-security-supported (PWG5100.22)
	"xri-security-supported": &DefAttr{
		SetOf: true,
		Min:   1,
		Max:   255,
		Tags:  []goipp.Tag{goipp.TagKeyword},
	},
	// System Status/xri-uri-scheme-supported (PWG5100.22)
	"xri-uri-scheme-supported": &DefAttr{
		SetOf: true,
		Min:   0,
		Max:   63,
		Tags:  []goipp.Tag{goipp.TagURIScheme},
	},
}

// Collections contains all top-level collections (groups) of
// attributes, indexed by name
var Collections = map[string]map[string]*DefAttr{
	"CUPS Device Attributes":        CUPSDeviceAttributes,
	"CUPS PPD Attributes":           CUPSPPDAttributes,
	"CUPS Printer Class Attributes": CUPSPrinterClassAttributes,
	"Document Description":          DocumentDescription,
	"Document Status":               DocumentStatus,
	"Document Template":             DocumentTemplate,
	"Event Notifications":           EventNotifications,
	"Job Description":               JobDescription,
	"Job Status":                    JobStatus,
	"Job Template":                  JobTemplate,
	"Operation":                     Operation,
	"Printer Description":           PrinterDescription,
	"Printer Status":                PrinterStatus,
	"Resource Description":          ResourceDescription,
	"Resource Status":               ResourceStatus,
	"Subscription Status":           SubscriptionStatus,
	"Subscription Template":         SubscriptionTemplate,
	"System Description":            SystemDescription,
	"System Status":                 SystemStatus,
}

// borrowings contains a table of attributes borrowing
// between collections.
var borrowings = []borrowing{
	{"Document Status/cover-back-actual", "Job Template/cover-back"},
	{"Document Status/cover-front-actual", "Job Template/cover-front"},
	{"Document Status/document-format-details", "Operation/document-format-details"},
	{"Document Status/document-format-details-detected", "Operation/document-format-details"},
	{"Document Status/finishings-col-actual", "Job Template/finishings-col"},
	{"Document Status/impressions-completed-col", "Document Status/impressions-col"},
	{"Document Status/input-attributes-actual", "Operation/input-attributes"},
	{"Document Status/insert-sheet-actual", "Job Template/insert-sheet"},
	{"Document Status/materials-col-actual", "Document Template/materials-col"},
	{"Document Status/media-col-actual", "Job Template/media-col"},
	{"Document Status/media-sheets-completed-col", "Document Status/media-sheets-col"},
	{"Document Status/output-attributes-actual", "Operation/output-attributes"},
	{"Document Status/pages-completed-col", "Document Status/pages-col"},
	{"Document Status/print-accuracy-actual", "Document Template/print-accuracy"},
	{"Document Status/print-objects-actual", "Document Template/print-objects"},
	{"Document Status/separator-sheets-actual", "Job Template/separator-sheets"},
	{"Document Template/cover-back", "Job Template/cover-back"},
	{"Document Template/cover-front", "Job Template/cover-front"},
	{"Document Template/finishings-col", "Job Template/finishings-col"},
	{"Document Template/insert-sheet", "Job Template/insert-sheet"},
	{"Document Template/media-col", "Job Template/media-col"},
	{"Document Template/separator-sheets", "Job Template/separator-sheets"},
	{"Job Status/client-info", "Operation/client-info"},
	{"Job Status/cover-back-actual", "Job Template/cover-back"},
	{"Job Status/cover-front-actual", "Job Template/cover-front"},
	{"Job Status/document-format-details-detected", "Operation/document-format-details"},
	{"Job Status/document-format-details-supplied", "Operation/document-format-details"},
	{"Job Status/finishings-col-actual", "Job Template/finishings-col"},
	{"Job Status/input-attributes-actual", "Operation/input-attributes"},
	{"Job Status/insert-sheet-actual", "Job Template/insert-sheet"},
	{"Job Status/job-accounting-sheets-actual", "Job Template/job-accounting-sheets"},
	{"Job Status/job-cover-back-actual", "Job Template/cover-back"},
	{"Job Status/job-cover-front-actual", "Job Template/cover-front"},
	{"Job Status/job-error-sheet-actual", "Job Template/job-error-sheet"},
	{"Job Status/job-impressions-completed-col", "Job Status/job-impressions-col"},
	{"Job Status/job-media-sheets-completed-col", "Job Status/job-media-sheets-col"},
	{"Job Status/job-pages-completed-col", "Job Status/job-pages-col"},
	{"Job Status/job-sheets-col-actual", "Job Template/job-sheets-col"},
	{"Job Status/job-storage", "Operation/job-storage"},
	{"Job Status/materials-col-actual", "Document Template/materials-col"},
	{"Job Status/media-col-actual", "Job Template/media-col"},
	{"Job Status/output-attributes-actual", "Operation/output-attributes"},
	{"Job Status/overrides-actual", "Job Template/overrides"},
	{"Job Status/print-accuracy-actual", "Document Template/print-accuracy"},
	{"Job Status/print-objects-actual", "Document Template/print-objects"},
	{"Job Status/separator-sheets-actual", "Job Template/separator-sheets"},
	{"Job Template/job-cover-back", "Job Template/cover-back"},
	{"Job Template/job-cover-front", "Job Template/cover-front"},
	{"Job Template/overrides", "Job Template"},
	{"Job Template/cover-back/media-col", "Job Template/media-col"},
	{"Job Template/cover-front/media-col", "Job Template/media-col"},
	{"Job Template/destination-uris/destination-attributes", "Document Template"},
	{"Job Template/destination-uris/destination-attributes", "Job Template"},
	{"Job Template/destination-uris/destination-attributes", "Operation"},
	{"Job Template/finishings-col/media-size", "Job Template/media-col/media-size"},
	{"Job Template/insert-sheet/media-col", "Job Template/media-col"},
	{"Job Template/job-accounting-sheets/media-col", "Job Template/media-col"},
	{"Job Template/job-error-sheet/media-col", "Job Template/media-col"},
	{"Job Template/job-sheets-col/media-col", "Job Template/media-col"},
	{"Job Template/proof-print/media-col", "Job Template/media-col"},
	{"Job Template/separator-sheets/media-col", "Job Template/media-col"},
	{"Operation/job-impressions-col", "Job Description/job-impressions-col"},
	{"Operation/job-media-sheets-col", "Job Description/job-media-sheets-col"},
	{"Operation/job-pages-col", "Job Status/job-pages-col"},
	{"Operation/preferred-attributes", "Document Template"},
	{"Operation/preferred-attributes", "Job Template"},
	{"Operation/preferred-attributes", "Operation"},
	{"Printer Description/cover-back-default", "Job Template/cover-back"},
	{"Printer Description/cover-front-default", "Job Template/cover-front"},
	{"Printer Description/cover-sheet-info-default", "Job Template/cover-sheet-info"},
	{"Printer Description/document-format-details-default", "Operation/document-format-details"},
	{"Printer Description/finishings-col-database", "Job Template/finishings-col"},
	{"Printer Description/finishings-col-default", "Job Template/finishings-col"},
	{"Printer Description/finishings-col-ready", "Job Template/finishings-col"},
	{"Printer Description/input-attributes-default", "Operation/input-attributes"},
	{"Printer Description/insert-sheet-default", "Job Template/insert-sheet"},
	{"Printer Description/job-accounting-sheets-default", "Job Template/job-accounting-sheets"},
	{"Printer Description/job-constraints-supported", "Job Template"},
	{"Printer Description/job-cover-back-default", "Job Template/cover-back"},
	{"Printer Description/job-cover-front-default", "Job Template/cover-front"},
	{"Printer Description/job-error-sheet-default", "Job Template/job-error-sheet"},
	{"Printer Description/job-presets-supported", "Job Template"},
	{"Printer Description/job-resolvers-supported", "Job Template"},
	{"Printer Description/job-sheets-col-default", "Job Template/job-sheets-col"},
	{"Printer Description/materials-col-database", "Document Template/materials-col"},
	{"Printer Description/materials-col-default", "Document Template/materials-col"},
	{"Printer Description/materials-col-ready", "Document Template/materials-col"},
	{"Printer Description/media-col-database", "Job Template/media-col"},
	{"Printer Description/media-col-default", "Job Template/media-col"},
	{"Printer Description/media-col-ready", "Job Template/media-col"},
	{"Printer Description/media-overprint-default", "Document Template/media-overprint"},
	{"Printer Description/output-attributes-default", "Operation/output-attributes"},
	{"Printer Description/proof-print-default", "Job Template/proof-print"},
	{"Printer Description/separator-sheets-default", "Job Template/separator-sheets"},
	{"Printer Description/destination-uri-ready/destination-attributes", "Document Template"},
	{"Printer Description/destination-uri-ready/destination-attributes", "Job Template"},
	{"Printer Description/destination-uri-ready/destination-attributes", "Operation"},
	{"Printer Status/printer-impressions-completed-col", "Job Description/job-impressions-col"},
	{"Printer Status/printer-media-sheets-completed-col", "Job Description/job-media-sheets-col"},
	{"Printer Status/printer-pages-completed-col", "Document Status/pages-col"},
	{"System Description/system-contact-col", "Printer Description/printer-contact-col"},
	{"System Description/system-service-contact-col", "Printer Description/printer-contact-col"},
	{"System Description/system-xri-supported", "Printer Description/printer-xri-supported"},
	{"System Status/system-impressions-completed-col", "Job Description/job-media-sheets-col"},
	{"System Status/system-media-sheets-completed-col", "Job Description/job-media-sheets-col"},
	{"System Status/system-pages-completed-col", "Document Status/pages-col"},
	{"System Status/system-configured-printers/printer-xri-supported", "Printer Description/printer-xri-supported"},
}

// exceptions contains member attributes that doesn't exist even if borrowed.
var exceptions = generic.NewSetOf(
	"Job Template/destination-uris/destination-attributes/document-password",
	"Job Template/destination-uris/destination-attributes/job-password",
	"Job Template/destination-uris/destination-attributes/job-password-encryption",
	"Printer Description/destination-uri-ready/destination-attributes/document-password",
	"Printer Description/destination-uri-ready/destination-attributes/job-password",
	"Printer Description/destination-uri-ready/destination-attributes/job-password-encryption",
)
