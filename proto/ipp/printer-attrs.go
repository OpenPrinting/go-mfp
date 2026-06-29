// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Printer Attributes

package ipp

import (
	"time"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// Commonly used Operation Attributes.
const (
	// AttrOperationAttributesCharset specifies character set
	// for all attributes in the message.
	//
	// It must be very first attribute in the message (RFC8011, 4.1.4.).
	AttrOperationAttributesCharset = "attributes-charset"

	// AttrOperationAttributesNaturalLanguage specifies natural
	// language for all textWithoutLanguage and nameWithoutLanguage
	// attributes in the message.
	//
	// It must be second attribute in the message (RFC8011, 4.1.4.).
	AttrOperationAttributesNaturalLanguage = "attributes-natural-language"

	// AttrOperationPrinterURI uses in requests to specify Printer URI,
	// where appropriate.
	AttrOperationPrinterURI = "printer-uri"

	// AttrOperationRequestedAttributes used in Get-Printer-Attributes
	// request to specify list of requested attributes.
	AttrOperationRequestedAttributes = "requested-attributes"

	// AttrOperationStatusMessage provides a short textual description
	// of the status of the operation.
	AttrOperationStatusMessage = "status-message"
)

// PrinterAttributes represents IPP Printer Attributes
type PrinterAttributes struct {
	ObjectRawAttrs

	PrinterDescriptionGroup
	PrinterStatusGroup

	PrinterDescription
	ScannerDescription
	JobTemplate
	MediaColDatabase
}

// PrinterDescription contains Printer Description and Status Attributes
type PrinterDescription struct {
	// RFC8011, 5.4: Printer Description and Status Attributes
	CharsetConfigured                 optional.Val[string]        `ipp:"charset-configured"`
	CharsetSupported                  []string                    `ipp:"charset-supported"`
	ColorSupported                    optional.Val[bool]          `ipp:"color-supported"`
	CompressionSupported              []KwCompression             `ipp:"compression-supported"`
	DocumentFormatDefault             optional.Val[string]        `ipp:"document-format-default"`
	DocumentFormatSupported           []string                    `ipp:"document-format-supported"`
	GeneratedNaturalLanguageSupported []string                    `ipp:"generated-natural-language-supported"`
	IppVersionsSupported              []goipp.Version             `ipp:"ipp-versions-supported"`
	JobImpressionsSupported           optional.Val[goipp.Range]   `ipp:"job-impressions-supported"`
	JobKOctetsSupported               optional.Val[goipp.Range]   `ipp:"job-k-octets-supported"`
	JobMediaSheetsSupported           optional.Val[goipp.Range]   `ipp:"job-media-sheets-supported"`
	MediaSizeSupported                []MediaSizeRange            `ipp:"media-size-supported"`
	MultipleDocumentJobsSupported     optional.Val[bool]          `ipp:"multiple-document-jobs-supported"`
	MultipleOperationTimeOut          optional.Val[int]           `ipp:"multiple-operation-time-out"`
	NaturalLanguageConfigured         optional.Val[string]        `ipp:"natural-language-configured"`
	OperationsSupported               []goipp.Op                  `ipp:"operations-supported"`
	PagesPerMinuteColor               optional.Val[int]           `ipp:"pages-per-minute-color"`
	PagesPerMinute                    optional.Val[int]           `ipp:"pages-per-minute"`
	PdlOverrideSupported              optional.Val[KwPdlOverride] `ipp:"pdl-override-supported"`
	PrinterDriverInstaller            optional.Val[string]        `ipp:"printer-driver-installer"`
	PrinterDeviceID                   optional.Val[string]        `ipp:"printer-device-id"`
	PrinterInfo                       optional.Val[string]        `ipp:"printer-info"`
	PrinterIsAcceptingJobs            optional.Val[bool]          `ipp:"printer-is-accepting-jobs"`
	PrinterLocation                   optional.Val[string]        `ipp:"printer-location"`
	PrinterMakeAndModel               optional.Val[string]        `ipp:"printer-make-and-model"`
	PrinterMessageFromOperator        optional.Val[string]        `ipp:"printer-message-from-operator"`
	PrinterMoreInfoManufacturer       optional.Val[string]        `ipp:"printer-more-info-manufacturer"`
	PrinterMoreInfo                   optional.Val[string]        `ipp:"printer-more-info"`
	PrinterName                       optional.Val[string]        `ipp:"printer-name"`
	PrinterStateMessage               optional.Val[string]        `ipp:"printer-state-message"`
	PrinterState                      optional.Val[int]           `ipp:"printer-state"`
	PrinterStateReasons               []KwPrinterStateReasons     `ipp:"printer-state-reasons"`
	PrinterUpTime                     optional.Val[int]           `ipp:"printer-up-time"`
	PrinterURISupported               []string                    `ipp:"printer-uri-supported"`
	QueuedJobCount                    optional.Val[int]           `ipp:"queued-job-count"`
	ReferenceURISchemesSupported      []string                    `ipp:"reference-uri-schemes-supported"`
	URIAuthenticationSupported        []KwURIAuthentication       `ipp:"uri-authentication-supported"`
	URISecuritySupported              []KwURISecurity             `ipp:"uri-security-supported"`

	// PWG5100.7: IPP Job Extensions v2.1 (JOBEXT)
	// 6.9 Printer Description Attributes
	ClientInfoSupported              []string                    `ipp:"client-info-supported"`
	DocumentCharsetDefault           optional.Val[string]        `ipp:"document-charset-default"`
	DocumentCharsetSupported         []string                    `ipp:"document-charset-supported"`
	DocumentFormatDetailsSupported   []string                    `ipp:"document-format-details-supported"`
	DocumentNaturalLanguageDefault   optional.Val[string]        `ipp:"document-natural-language-default"`
	DocumentNaturalLanguageSupported []string                    `ipp:"document-natural-language-supported"`
	JobCreationAttributesSupported   []string                    `ipp:"job-creation-attributes-supported"`
	JobHistoryAttributesConfigured   []string                    `ipp:"job-history-attributes-configured"`
	JobHistoryAttributesSupported    []string                    `ipp:"job-history-attributes-supported"`
	JobHistoryIntervalConfigured     optional.Val[int]           `ipp:"job-history-interval-configured"`
	JobHistoryIntervalSupported      optional.Val[goipp.Range]   `ipp:"job-history-interval-supported"`
	JobIdsSupported                  optional.Val[bool]          `ipp:"job-ids-supported"`
	JobMandatoryAttributesSupported  optional.Val[bool]          `ipp:"job-mandatory-attributes-supported"`
	JobSpoolingSupported             optional.Val[KwJobSpooling] `ipp:"job-spooling-supported"`
	MediaBackCoatingSupported        []KwMediaBackCoating        `ipp:"media-back-coating-supported"`
	MediaBottomMarginSupported       []int                       `ipp:"media-bottom-margin-supported"`
	MediaColDefault                  optional.Val[MediaCol]      `ipp:"media-col-default"`
	MediaColorSupported              []string                    `ipp:"media-color-supported"`
	MediaColReady                    []MediaColEx                `ipp:"media-col-ready"`
	MediaColSupported                []string                    `ipp:"media-col-supported"`
	MediaFrontCoatingSupported       []KwMediaBackCoating        `ipp:"media-front-coating-supported"`
	MediaGrainSupported              []string                    `ipp:"media-grain-supported"`
	MediaHoleCountSupported          []goipp.Range               `ipp:"media-hole-count-supported"`
	MediaKeySupported                []KwMedia                   `ipp:"media-key-supported"`
	MediaLeftMarginSupported         []int                       `ipp:"media-left-margin-supported"`
	MediaOrderCountSupported         []goipp.Range               `ipp:"media-order-count-supported"`
	MediaPrePrintedSupported         []string                    `ipp:"media-pre-printed-supported"`
	MediaRecycledSupported           []string                    `ipp:"media-recycled-supported"`
	MediaRightMarginSupported        []int                       `ipp:"media-right-margin-supported"`
	MediaSourceSupported             []string                    `ipp:"media-source-supported"`
	MediaThicknessSupported          []goipp.Range               `ipp:"media-thickness-supported"`
	MediaToothSupported              []string                    `ipp:"media-tooth-supported"`
	MediaTopMarginSupported          []int                       `ipp:"media-top-margin-supported"`
	MediaTypeSupported               []string                    `ipp:"media-type-supported"`
	MediaWeightMetricSupported       []goipp.Range               `ipp:"media-weight-metric-supported"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.5 Printer Description Attributes
	DocumentPasswordSupported         optional.Val[int]         `ipp:"document-password-supported"`
	IdentifyActionsDefault            []string                  `ipp:"identify-actions-default"`
	IdentifyActionsSupported          []string                  `ipp:"identify-actions-supported"`
	// Scan Services MUST include "scan" in this list (PWG5100.17, 9.2).
	IppFeaturesSupported              []string                  `ipp:"ipp-features-supported"`
	JobPresetsSupported               []JobPresets              `ipp:"job-presets-supported"`
	JpegFeaturesSupported             []string                  `ipp:"jpeg-features-supported"`
	JpegKOctetsSupported              optional.Val[goipp.Range] `ipp:"jpeg-k-octets-supported"`
	JpegXDimensionSupported           optional.Val[goipp.Range] `ipp:"jpeg-x-dimension-supported"`
	JpegYDimensionSupported           optional.Val[goipp.Range] `ipp:"jpeg-y-dimension-supported"`
	MultipleOperationTimeOutAction    optional.Val[string]      `ipp:"multiple-operation-time-out-action"`
	PdfKOctetsSupported               optional.Val[goipp.Range] `ipp:"pdf-k-octets-supported"`
	PdfVersionsSupported              []string                  `ipp:"pdf-versions-supported"`
	PreferredAttributesSupported      optional.Val[bool]        `ipp:"preferred-attributes-supported"`
	PrinterDNSSdName                  optional.Val[string]      `ipp:"printer-dns-sd-name"`
	PrinterGeoLocation                optional.Val[string]      `ipp:"printer-geo-location"`
	PrinterGetAttributesSupported     []string                  `ipp:"printer-get-attributes-supported"`
	PrinterIcons                      []string                  `ipp:"printer-icons"`
	PrinterKind                       []string                  `ipp:"printer-kind"`
	PrinterOrganization               []string                  `ipp:"printer-organization"`
	PrinterOrganizationalUnit         []string                  `ipp:"printer-organizational-unit"`
	PrinterStringsLanguagesSupported  []string                  `ipp:"printer-strings-languages-supported"`
	PrinterStringsURI                 optional.Val[string]      `ipp:"printer-strings-uri"`
	RequestingUserURISupported        optional.Val[bool]        `ipp:"requesting-user-uri-supported"`
	RequestingUserURISchemesSupported []string                  `ipp:"requesting-user-uri-schemes-supported"`

	// PWG5100.13: IPP Driver Replacement Extensions v2.0 (NODRIVER)
	// 6.6 Printer Status Attributes
	DeviceServiceCount           optional.Val[int]       `ipp:"device-service-count"`
	DeviceUUID                   optional.Val[string]    `ipp:"device-uuid"`
	PrinterConfigChangeDateTime  optional.Val[time.Time] `ipp:"printer-config-change-date-time"`
	PrinterConfigChangeTime      optional.Val[int]       `ipp:"printer-config-change-time"`
	PrinterFirmwareName          []string                `ipp:"printer-firmware-name"`
	PrinterFirmwarePatches       []string                `ipp:"printer-firmware-patches"`
	PrinterFirmwareStringVersion []string                `ipp:"printer-firmware-string-version"`
	PrinterFirmwareVersion       []string                `ipp:"printer-firmware-version"`
	PrinterInputTray             []string                `ipp:"printer-input-tray"`
	PrinterOutputTray            []string                `ipp:"printer-output-tray"`
	PrinterSupplyDescription     []goipp.TextWithLang    `ipp:"printer-supply-description"`
	PrinterSupplyInfoURI         optional.Val[string]    `ipp:"printer-supply-info-uri"`
	PrinterSupply                []string                `ipp:"printer-supply"`
	PrinterUUID                  optional.Val[string]    `ipp:"printer-uuid"`

	// Wi-Fi Peer-to-Peer Services Print (P2Ps-Print)
	// Technical Specification
	// (for Wi-Fi Direct® services certification)
	PclmRasterBackSide       optional.Val[string] `ipp:"pclm-raster-back-side"`
	PclmStripHeightPreferred optional.Val[int]    `ipp:"pclm-strip-height-preferred"`
	PclmStripHeightSupported []int                `ipp:"pclm-strip-height-supported"`

	// CUPS extensions
	DeviceURI          string                      `ipp:"device-uri"`
	MarkerChangeTime   optional.Val[int]           `ipp:"marker-change-time"`
	MarkerColors       []string                    `ipp:"marker-colors"`
	MarkerHighLevels   []int                       `ipp:"marker-high-levels"`
	MarkerLevels       []int                       `ipp:"marker-levels"`
	MarkerLowLevels    []int                       `ipp:"marker-low-levels"`
	MarkerMessage      optional.Val[string]        `ipp:"marker-message"`
	MarkerNames        []string                    `ipp:"marker-names"`
	MarkerTypes        []string                    `ipp:"marker-types"`
	PrinterID          optional.Val[int]           `ipp:"printer-id"`
	PrinterIsShared    optional.Val[bool]          `ipp:"printer-is-shared"`
	PrinterIsTemporary optional.Val[bool]          `ipp:"printer-is-temporary"`
	PrinterType        optional.Val[EnPrinterType] `ipp:"printer-type"`
	UrfSupported       []string                    `ipp:"urf-supported"`
}

// PrinterJobSaveDisposition represents "job-save-disposition-default"
// collection entry in PrinterAttributes
type PrinterJobSaveDisposition struct {
	SaveDisposition optional.Val[string] `ipp:"save-disposition"`
	SaveInfo        []PrinterSaveInfo    `ipp:"save-info"`
}

// PrinterSaveInfo represents "save-info" collection entry
// in PrinterJobSaveDisposition
type PrinterSaveInfo struct {
	SaveLocation       optional.Val[string] `ipp:"save-location"`
	SaveName           optional.Val[string] `ipp:"save-name"`
	SaveDocumentFormat optional.Val[string] `ipp:"save-document-format"`
}

// DecodePrinterAttributes decodes [PrinterAttributes] from
// [goipp.Attributes].
func DecodePrinterAttributes(attrs goipp.Attributes, opt *DecoderOptions) (
	*PrinterAttributes, error) {

	pa := &PrinterAttributes{}
	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(pa, attrs)
	if err != nil {
		return nil, err
	}
	return pa, nil
}

// IsCharsetSupported tells if charset is supported
func (pa *PrinterAttributes) IsCharsetSupported(cs string) bool {
	for _, supp := range pa.CharsetSupported {
		if cs == supp {
			return true
		}
	}
	return false
}

// IsOperationSupported tells if operation is supported
func (pa *PrinterAttributes) IsOperationSupported(op goipp.Op) bool {
	for _, supp := range pa.OperationsSupported {
		if op == supp {
			return true
		}
	}
	return false
}
