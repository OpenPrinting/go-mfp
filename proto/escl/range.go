// MFP - Miulti-Function Printers and scanners toolkit
// eSCL core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common type for Range of some value.

package escl

import (
	"strconv"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Range commonly used to specify the range of some parameter, like
// brightness, contrast etc.
type Range struct {
	Min    int               // Minimal supported value
	Max    int               // Maximal supported value
	Normal int               // Normal value
	Step   optional.Val[int] // Step between the subsequent values
}

// decodeRange decodes [Range] from the XML tree
func decodeRange(root xmldoc.Element) (r Range, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	// The Range element is a special case. The Mopria eSCL specification
	// references it as scan:RangeAndStepOfIntType. However, it fails to
	// define the type itself, and the official examples lack any usage
	// cases.
	//
	// As a result, some firmwares use invalid namespace prefixes for the
	// child elements, or omit namespace prefixes entirely.
	//
	// For example, the Xerox B235 encodes this element as follows:
	//
	//    <scan:BrightnessSupport>
	//       <Min>1</Min>
	//       <Max>9</Max>
	//       <Normal>5</Normal>
	//       <Step>1</Step>
	//    </scan:BrightnessSupport>
	//
	// This causes our strict XML parser to fail due to missing expected
	// elements (e.g., <scan:Min>).
	//
	// As a workaround, we force all child elements of a Range block
	// to use the "scan:" prefix.
	//
	// This solution is not perfect, as it silently bypasses a firmware
	// bug without logging a warning or preserving information about the
	// original prefix.
	//
	// But for now, this is better than nothing. -- FIXME
	for i := range root.Children {
		chld := &root.Children[i]
		if i := strings.IndexByte(chld.Name, ':'); i >= 0 {
			chld.Name = NsScan + chld.Name[i:]
		} else {
			chld.Name = NsScan + ":" + chld.Name
		}
	}

	// Lookup relevant XML elements
	min := xmldoc.Lookup{Name: NsScan + ":Min", Required: true}
	max := xmldoc.Lookup{Name: NsScan + ":Max", Required: true}
	normal := xmldoc.Lookup{Name: NsScan + ":Normal", Required: true}
	step := xmldoc.Lookup{Name: NsScan + ":Step"}

	missed := root.Lookup(&min, &max, &normal, &step)
	if missed != nil {
		err = xmldoc.XMLErrMissed(missed.Name)
		return
	}

	// Decode elements
	r.Min, err = decodeInt(min.Elem)
	if err == nil {
		r.Max, err = decodeInt(max.Elem)
	}
	if err == nil {
		r.Normal, err = decodeInt(normal.Elem)
	}
	if err == nil && step.Found {
		var tmp int
		tmp, err = decodeNonNegativeInt(step.Elem)
		r.Step = optional.New(tmp)
	}

	return
}

// toXML generates XML tree for the [Range].
func (r Range) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			{
				Name: NsScan + ":" + "Min",
				Text: strconv.FormatInt(int64(r.Min), 10),
			},
			{
				Name: NsScan + ":" + "Max",
				Text: strconv.FormatInt(int64(r.Max), 10),
			},
			{
				Name: NsScan + ":" + "Normal",
				Text: strconv.FormatInt(int64(r.Normal), 10),
			},
		},
	}

	if r.Step != nil {
		step := xmldoc.Element{
			Name: NsScan + ":" + "Step",
			Text: strconv.FormatInt(int64(*r.Step), 10),
		}
		elm.Children = append(elm.Children, step)
	}

	return elm
}
