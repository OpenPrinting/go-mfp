// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// device settings

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// DeviceSettings represents the <wscn:DeviceSettings> element,
// describing the basic capabilities of the scan device.
type DeviceSettings struct {
	AutoExposureSupported             BooleanElement
	BrightnessSupported               BooleanElement
	CompressionQualityFactorSupported Range
	ContentTypesSupported             []ContentTypeValue
	ContrastSupported                 BooleanElement
	DocumentSizeAutoDetectSupported   BooleanElement
	FormatsSupported                  []FormatValue
	RotationsSupported                []RotationValue
	ScalingRangeSupported             ScalingRangeSupported
}

// toXML generates XML tree for the [DeviceSettings].
func (ds DeviceSettings) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		ds.AutoExposureSupported.toXML(NsWSCN + ":AutoExposureSupported"),
		ds.BrightnessSupported.toXML(NsWSCN + ":BrightnessSupported"),
		{
			Name:     NsWSCN + ":CompressionQualityFactorSupported",
			Children: ds.CompressionQualityFactorSupported.toXML(),
		},
	}
	// ContentTypesSupported
	ctChildren := make([]xmldoc.Element, len(ds.ContentTypesSupported))
	for i, v := range ds.ContentTypesSupported {
		ctChildren[i] = v.toXML(NsWSCN + ":ContentTypeValue")
	}
	children = append(children, xmldoc.Element{
		Name:     NsWSCN + ":ContentTypesSupported",
		Children: ctChildren,
	})
	children = append(children,
		ds.ContrastSupported.toXML(NsWSCN+":ContrastSupported"),
		ds.DocumentSizeAutoDetectSupported.toXML(NsWSCN+
			":DocumentSizeAutoDetectSupported"),
	)
	// FormatsSupported
	fmtChildren := make([]xmldoc.Element, len(ds.FormatsSupported))
	for i, v := range ds.FormatsSupported {
		fmtChildren[i] = v.toXML(NsWSCN + ":FormatValue")
	}
	children = append(children, xmldoc.Element{
		Name:     NsWSCN + ":FormatsSupported",
		Children: fmtChildren,
	})
	// RotationsSupported
	rotChildren := make([]xmldoc.Element, len(ds.RotationsSupported))
	for i, v := range ds.RotationsSupported {
		rotChildren[i] = v.toXML(NsWSCN + ":RotationValue")
	}
	children = append(children, xmldoc.Element{
		Name:     NsWSCN + ":RotationsSupported",
		Children: rotChildren,
	})
	children = append(children, ds.ScalingRangeSupported.toXML(NsWSCN+
		":ScalingRangeSupported"))
	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeDeviceSettings decodes [DeviceSettings] from the XML tree.
func decodeDeviceSettings(root xmldoc.Element) (ds DeviceSettings, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	// Lookup all required XML elements
	autoExposureSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":AutoExposureSupported",
		Required: true,
	}
	brightnessSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":BrightnessSupported",
		Required: true,
	}
	compressionQualityFactorSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":CompressionQualityFactorSupported",
		Required: true,
	}
	contentTypesSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":ContentTypesSupported",
		Required: true,
	}
	contrastSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":ContrastSupported",
		Required: true,
	}
	documentSizeAutoDetectSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":DocumentSizeAutoDetectSupported",
		Required: true,
	}
	formatsSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":FormatsSupported",
		Required: true,
	}
	rotationsSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":RotationsSupported",
		Required: true,
	}
	scalingRangeSupported := xmldoc.Lookup{
		Name:     NsWSCN + ":ScalingRangeSupported",
		Required: true,
	}

	missed := root.Lookup(
		&autoExposureSupported,
		&brightnessSupported,
		&compressionQualityFactorSupported,
		&contentTypesSupported,
		&contrastSupported,
		&documentSizeAutoDetectSupported,
		&formatsSupported,
		&rotationsSupported,
		&scalingRangeSupported,
	)
	if missed != nil {
		return ds, xmldoc.XMLErrMissed(missed.Name)
	}

	if ds.AutoExposureSupported, err = decodeBooleanElement(
		autoExposureSupported.Elem); err != nil {
		return ds, fmt.Errorf("AutoExposureSupported: %w", err)
	}
	if ds.BrightnessSupported, err = decodeBooleanElement(
		brightnessSupported.Elem); err != nil {
		return ds, fmt.Errorf("BrightnessSupported: %w", err)
	}
	if ds.CompressionQualityFactorSupported, err = decodeRange(compressionQualityFactorSupported.Elem); err != nil {
		return ds, fmt.Errorf("CompressionQualityFactorSupported: %w", err)
	}

	// ContentTypesSupported slice
	for _, child := range contentTypesSupported.Elem.Children {
		if child.Name == NsWSCN+":ContentTypeValue" {
			val, err := decodeContentTypeValue(child)
			if err != nil {
				return ds, fmt.Errorf("ContentTypesSupported: "+
					"invalid ContentTypeValue: %w", err)
			}
			ds.ContentTypesSupported = append(ds.ContentTypesSupported, val)
		}
	}
	if len(ds.ContentTypesSupported) == 0 {
		return ds, fmt.Errorf("ContentTypesSupported: " +
			"at least one ContentTypeValue is required")
	}
	if ds.ContrastSupported, err = decodeBooleanElement(
		contrastSupported.Elem); err != nil {
		return ds, fmt.Errorf("ContrastSupported: %w", err)
	}
	if ds.DocumentSizeAutoDetectSupported, err = decodeBooleanElement(documentSizeAutoDetectSupported.Elem); err != nil {
		return ds, fmt.Errorf("DocumentSizeAutoDetectSupported: %w", err)
	}
	// FormatsSupported slice
	for _, child := range formatsSupported.Elem.Children {
		if child.Name == NsWSCN+":FormatValue" {
			val, err := decodeFormatValue(child)
			if err != nil {
				return ds, fmt.Errorf("FormatsSupported: "+
					"invalid FormatValue: %w", err)
			}
			ds.FormatsSupported = append(ds.FormatsSupported, val)
		}
	}
	if len(ds.FormatsSupported) == 0 {
		return ds, fmt.Errorf("FormatsSupported: " +
			"at least one FormatValue is required")
	}
	// RotationsSupported slice
	for _, child := range rotationsSupported.Elem.Children {
		if child.Name == NsWSCN+":RotationValue" {
			val, err := decodeRotationValue(child)
			if err != nil {
				return ds, fmt.Errorf("RotationsSupported: "+
					"invalid RotationValue: %w", err)
			}
			ds.RotationsSupported = append(ds.RotationsSupported, val)
		}
	}
	if len(ds.RotationsSupported) == 0 {
		return ds, fmt.Errorf("RotationsSupported: " +
			"at least one RotationValue is required")
	}
	if ds.ScalingRangeSupported, err = decodeScalingRangeSupported(
		scalingRangeSupported.Elem,
	); err != nil {
		return ds, fmt.Errorf("ScalingRangeSupported: %w", err)
	}
	return ds, nil
}
