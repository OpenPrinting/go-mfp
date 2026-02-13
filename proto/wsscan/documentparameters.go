// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// DocumentParameters:
// defines image processing functions to be applied to documents

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// DocumentParameters defines the image processing functions to be applied
// to documents within a scan job. All child elements are optional.
type DocumentParameters struct {
	CompressionQualityFactor optional.Val[CompressionQualityFactor]
	ContentType              optional.Val[ContentType]
	Exposure                 optional.Val[Exposure]
	FilmScanMode             optional.Val[FilmScanModeElement]
	Format                   optional.Val[Format]
	ImagesToTransfer         optional.Val[ImagesToTransfer]
	InputSize                optional.Val[InputSize]
	InputSource              optional.Val[InputSource]
	MediaSides               optional.Val[MediaSides]
	Rotation                 optional.Val[Rotation]
	Scaling                  optional.Val[Scaling]
}

// toXML generates XML tree for the DocumentParameters.
func (dp DocumentParameters) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{}

	if dp.CompressionQualityFactor != nil {
		children = append(children, optional.Get(
			dp.CompressionQualityFactor).toXML(
			NsWSCN+":CompressionQualityFactor"))
	}

	if dp.ContentType != nil {
		children = append(children, optional.Get(
			dp.ContentType).toXML(
			NsWSCN+":ContentType"))
	}

	if dp.Exposure != nil {
		children = append(children, optional.Get(
			dp.Exposure).toXML(
			NsWSCN+":Exposure"))
	}

	if dp.FilmScanMode != nil {
		children = append(children, optional.Get(
			dp.FilmScanMode).toXML(
			NsWSCN+":FilmScanMode"))
	}

	if dp.Format != nil {
		children = append(children, optional.Get(
			dp.Format).toXML(
			NsWSCN+":Format"))
	}

	if dp.ImagesToTransfer != nil {
		children = append(children, optional.Get(
			dp.ImagesToTransfer).toXML(
			NsWSCN+":ImagesToTransfer"))
	}

	if dp.InputSize != nil {
		children = append(children, optional.Get(
			dp.InputSize).toXML(
			NsWSCN+":InputSize"))
	}

	if dp.InputSource != nil {
		children = append(children, optional.Get(
			dp.InputSource).toXML(
			NsWSCN+":InputSource"))
	}

	if dp.MediaSides != nil {
		children = append(children, optional.Get(
			dp.MediaSides).toXML(
			NsWSCN+":MediaSides"))
	}

	if dp.Rotation != nil {
		children = append(children, optional.Get(
			dp.Rotation).toXML(
			NsWSCN+":Rotation"))
	}

	if dp.Scaling != nil {
		children = append(children, optional.Get(
			dp.Scaling).toXML(
			NsWSCN+":Scaling"))
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeDocumentParameters decodes DocumentParameters from the XML tree.
func decodeDocumentParameters(root xmldoc.Element) (
	dp DocumentParameters,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	compressionQualityFactor := xmldoc.Lookup{
		Name:     NsWSCN + ":CompressionQualityFactor",
		Required: false,
	}
	contentType := xmldoc.Lookup{
		Name:     NsWSCN + ":ContentType",
		Required: false,
	}
	exposure := xmldoc.Lookup{
		Name:     NsWSCN + ":Exposure",
		Required: false,
	}
	filmScanMode := xmldoc.Lookup{
		Name:     NsWSCN + ":FilmScanMode",
		Required: false,
	}
	format := xmldoc.Lookup{
		Name:     NsWSCN + ":Format",
		Required: false,
	}
	imagesToTransfer := xmldoc.Lookup{
		Name:     NsWSCN + ":ImagesToTransfer",
		Required: false,
	}
	inputSize := xmldoc.Lookup{
		Name:     NsWSCN + ":InputSize",
		Required: false,
	}
	inputSource := xmldoc.Lookup{
		Name:     NsWSCN + ":InputSource",
		Required: false,
	}
	mediaSides := xmldoc.Lookup{
		Name:     NsWSCN + ":MediaSides",
		Required: false,
	}
	rotation := xmldoc.Lookup{
		Name:     NsWSCN + ":Rotation",
		Required: false,
	}
	scaling := xmldoc.Lookup{
		Name:     NsWSCN + ":Scaling",
		Required: false,
	}

	root.Lookup(
		&compressionQualityFactor,
		&contentType,
		&exposure,
		&filmScanMode,
		&format,
		&imagesToTransfer,
		&inputSize,
		&inputSource,
		&mediaSides,
		&rotation,
		&scaling,
	)

	// Decode each optional element if present
	if compressionQualityFactor.Found {
		var cqf CompressionQualityFactor
		if cqf, err = decodeCompressionQualityFactor(
			compressionQualityFactor.Elem,
		); err != nil {
			return dp, fmt.Errorf("CompressionQualityFactor: %w", err)
		}
		dp.CompressionQualityFactor = optional.New(cqf)
	}

	if contentType.Found {
		var ct ContentType
		if ct, err = decodeContentType(
			contentType.Elem,
		); err != nil {
			return dp, fmt.Errorf("ContentType: %w", err)
		}
		dp.ContentType = optional.New(ct)
	}

	if exposure.Found {
		var exp Exposure
		if exp, err = decodeExposure(
			exposure.Elem,
		); err != nil {
			return dp, fmt.Errorf("Exposure: %w", err)
		}
		dp.Exposure = optional.New(exp)
	}

	if filmScanMode.Found {
		var fsm FilmScanModeElement
		if fsm, err = decodeFilmScanModeElement(
			filmScanMode.Elem,
		); err != nil {
			return dp, fmt.Errorf("FilmScanMode: %w", err)
		}
		dp.FilmScanMode = optional.New(fsm)
	}

	if format.Found {
		var fmtVal Format
		if fmtVal, err = decodeFormat(
			format.Elem,
		); err != nil {
			return dp, fmt.Errorf("Format: %w", err)
		}
		dp.Format = optional.New(fmtVal)
	}

	if imagesToTransfer.Found {
		var itt ImagesToTransfer
		if itt, err = decodeImagesToTransfer(
			imagesToTransfer.Elem,
		); err != nil {
			return dp, fmt.Errorf("ImagesToTransfer: %w", err)
		}
		dp.ImagesToTransfer = optional.New(itt)
	}

	if inputSize.Found {
		var is InputSize
		if is, err = decodeInputSize(
			inputSize.Elem,
		); err != nil {
			return dp, fmt.Errorf("InputSize: %w", err)
		}
		dp.InputSize = optional.New(is)
	}

	if inputSource.Found {
		var isrc InputSource
		if isrc, err = decodeInputSource(
			inputSource.Elem,
		); err != nil {
			return dp, fmt.Errorf("InputSource: %w", err)
		}
		dp.InputSource = optional.New(isrc)
	}

	if mediaSides.Found {
		var ms MediaSides
		if ms, err = decodeMediaSides(
			mediaSides.Elem,
		); err != nil {
			return dp, fmt.Errorf("MediaSides: %w", err)
		}
		dp.MediaSides = optional.New(ms)
	}

	if rotation.Found {
		var rot Rotation
		if rot, err = decodeRotation(
			rotation.Elem,
		); err != nil {
			return dp, fmt.Errorf("Rotation: %w", err)
		}
		dp.Rotation = optional.New(rot)
	}

	if scaling.Found {
		var scl Scaling
		if scl, err = decodeScaling(
			scaling.Elem,
		); err != nil {
			return dp, fmt.Errorf("Scaling: %w", err)
		}
		dp.Scaling = optional.New(scl)
	}

	return dp, nil
}
