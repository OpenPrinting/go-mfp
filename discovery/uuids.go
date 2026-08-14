// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Functions for management collections of UUIDs

package discovery

import (
	"sort"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// uuidsContain reports if collection of UUIDs contains
// the specified UUID.
//
// UUIDs assumed to be sorted.
func uuidsContain(uuids []uuid.UUID, uu uuid.UUID) bool {
	_, found := sort.Find(len(uuids), func(i int) int {
		return uu.Compare(uuids[i])
	})
	return found
}

// uuidsAdd adds UUID into the collection of UUIDs.
//
// It assumes sorted UUIDs on input and returns updated and
// properly sorted collection on output. It works by modifying
// the input collection in-place, so be aware.
//
// Additionally, it returns a bool flag that indicates that
// UUID was actually added.
func uuidsAdd(uuids []uuid.UUID, uu uuid.UUID) ([]uuid.UUID, bool) {
	i, found := sort.Find(len(uuids), func(i int) int {
		return uu.Compare(uuids[i])
	})

	if found {
		return uuids, false
	}

	uuids = append(uuids, uuid.UUID{})
	copy(uuids[i+1:], uuids[i:])
	uuids[i] = uu

	return uuids, true
}

// uuidsDel deletes UUID from the collection of UUIDs.
//
// It assumes sorted UUIDs on input and returns updated and
// properly sorted collection on output. It works by modifying
// the input collection in-place, so be aware.
//
// Additionally, it returns a bool flag that indicates that
// UUID was actually found in the collection and deleted.
func uuidsDel(uuids []uuid.UUID, uu uuid.UUID) ([]uuid.UUID, bool) {
	i, found := sort.Find(len(uuids), func(i int) int {
		return uu.Compare(uuids[i])
	})

	if !found {
		return uuids, false
	}

	copy(uuids[i:], uuids[i+1:])
	uuids = uuids[:len(uuids)-1]

	return uuids, true
}

// uuidsMerge merges two collections of UUIDs.
//
// Input collections must be sorted and returned collection
// is sorted as well.
func uuidsMerge(uuids1, uuids2 []uuid.UUID) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(uuids1)+len(uuids2))

	for len(uuids1) > 0 && len(uuids2) > 0 {
		cmp := uuids1[0].Compare(uuids2[0])
		switch {
		case cmp < 0:
			out = append(out, uuids1[0])
			uuids1 = uuids1[1:]
		case cmp > 0:
			out = append(out, uuids2[0])
			uuids2 = uuids2[1:]
		default:
			out = append(out, uuids1[0])
			uuids1 = uuids1[1:]
			uuids2 = uuids2[1:]
		}
	}

	out = append(out, uuids1...)
	out = append(out, uuids2...)

	return out
}

// uuidsOverlap reports if two collections of UUIDs overlap,
// i.e., contains some UUIDs in common.
//
// Input collections must be sorted.
func uuidsOverlap(uuids1, uuids2 []uuid.UUID) bool {
	for len(uuids1) > 0 && len(uuids2) > 0 {
		cmp := uuids1[0].Compare(uuids2[0])
		switch {
		case cmp < 0:
			uuids1 = uuids1[1:]
		case cmp > 0:
			uuids2 = uuids2[1:]
		default:
			return true
		}
	}

	return false
}
