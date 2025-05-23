// MFP - Miulti-Function Printers and scanners toolkit
// XML mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// XML elements lookup

package xmldoc

// Lookup contains element name for lookup by name and
// received lookup result.
//
// It is optimized for looking up multiple elements at once.
//
// See also: [Element.Lookup]
type Lookup struct {
	Name     string  // Requested element name
	Required bool    // This is required element
	Elem     Element // Returned element data
	Found    bool    // Becomes true, if element was found
}

// LookupAttr is like [Lookup], but for attributes, not for children.
type LookupAttr struct {
	Name     string // Requested element name
	Required bool   // This is required element
	Attr     Attr   // Returned element data
	Found    bool   // Becomes true, if element was found
}
