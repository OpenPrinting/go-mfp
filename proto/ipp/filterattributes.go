// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Attributes filtering for Get-XXX-Attributes

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/goipp"
)

// filterAttributes performs attribute filtering for
// Get-XXX-Attributes operations
type filterAttributes struct {
	// groups maps known groups of attributes (e.g., "all")
	// into concrete sets of attribute names
	groups map[string]generic.Set[string]

	// standard contains set of known standard (per IANA registrtion)
	// attribute names
	standard generic.Set[string]
}

// Apply filters encoded attributes against the requested groups/names
// defined by attrGroups. Returns the filtered subset (preserving the
// original order of encoded) and the list of requested names that were
// neither a known group nor a supported attribute.
func (filter *filterAttributes) Apply(
	requestedAttrs []string,
	encoded goipp.Attributes,
) (filtered goipp.Attributes, unsupported []string) {

	// Roll over all attributes, supported by device. Build
	// the following sets:
	//   - supported -- all attributes, supported by device
	//   - nonstandard -- attributes, supported by device,
	//     but not found in the standard groups
	supported := generic.NewSet[string]()
	nonstandard := generic.NewSet[string]()

	for _, attr := range encoded {
		supported.Add(attr.Name)
		if !filter.standard.Contains(attr.Name) {
			nonstandard.Add(attr.Name)
		}
	}

	// Roll over the requested attributes. Build the following sets:
	//   - expanded -- set of requested attributes with groups
	//     expanded and unknown attributes excluded
	//   - seenUnsupported is the set of requested but unsupported
	//     attributes
	expanded := generic.NewSet[string]()
	seenUnsupported := generic.NewSet[string]()

	for _, name := range requestedAttrs {
		if group, ok := filter.groups[name]; ok {
			expanded.Merge(group)
			if name == "all" {
				expanded.Merge(nonstandard)
			}
		} else if supported.Contains(name) {
			expanded.Add(name)
		} else if seenUnsupported.TestAndAdd(name) {
			unsupported = append(unsupported, name)
		}
	}

	// Now match present attributes against the expanded
	for _, attr := range encoded {
		if expanded.Contains(attr.Name) {
			filtered = append(filtered, attr)
		}
	}

	return
}
