// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// FilmScanMode element (not to be confused with FilmScanMode enum type)

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// FilmScanModeElement represents the optional <wscn:FilmScanMode> element
// that specifies the exposure type of the film to be scanned.
//
// Standard values include: "NotApplicable", "ColorSlideFilm",
// "ColorNegativeFilm", "BlackandWhiteNegativeFilm".
// Values can be extended or subset, so any string value is accepted.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
//
// Note: This is different from the [FilmScanMode] enum type which is used
// in FilmScanModesSupported lists.
type FilmScanModeElement = AttributedElement[string]

// decodeFilmScanModeElement decodes [FilmScanModeElement] from the XML tree.
func decodeFilmScanModeElement(root xmldoc.Element) (FilmScanModeElement, error) {
	return decodeAttributedElement(root, func(s string) (string, error) {
		// Accept any string value as values can be extended/subset
		return s, nil
	})
}

// toXMLFilmScanModeElement generates XML tree for the [FilmScanModeElement].
func toXMLFilmScanModeElement(fsm FilmScanModeElement, name string) xmldoc.Element {
	return fsm.toXML(name, func(s string) string {
		return s
	})
}
