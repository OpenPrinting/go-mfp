// MFP - Multi-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2026 and up by Anand Biju
// See LICENSE for license terms and conditions
//
// Output generation tests

package discovery

import (
	"context"
	"net/netip"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// countUnits returns total number of units in all groups
func countUnits(groups []unitgroup) int {
	n := 0
	for _, grp := range groups {
		n += len(grp)
	}
	return n
}

// TestMergeUnitGroupsByUUIDKeepsGroupsWithoutUUID checks that
// mergeUnitGrousByUUID passes groups without UUID through,
// instead of dropping them.
func TestMergeUnitGroupsByUUIDKeepsGroupsWithoutUUID(t *testing.T) {
	out := output{ctx: context.Background()}
	addr := []netip.Addr{netip.MustParseAddr("192.168.0.10")}
	uu := uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef")

	tests := []struct {
		name   string
		groups []unitgroup
		want   int
	}{
		{
			name: "all groups without UUID",
			groups: []unitgroup{
				{{ID: UnitID{Queue: "lpd1"}}},
				{{ID: UnitID{Queue: "lpd2"}}},
			},
			want: 2,
		},
		{
			name: "group without UUID and group with UUID",
			groups: []unitgroup{
				{{ID: UnitID{Queue: "lpd"}}},
				{{ID: UnitID{DNSSDName: "ipp", UUID: uu}}},
			},
			want: 2,
		},
		{
			name: "groups with the same UUID are still merged",
			groups: []unitgroup{
				{{ID: UnitID{Queue: "lpd"}}},
				// Backends must differ for the merge, so
				// use mockBackend and *mockBackend
				{{
					ID: UnitID{
						DNSSDName: "ipp",
						UUID:      uu,
						Backend:   &mockBackend{},
					},
					Addrs: addr,
				}},
				{{
					ID: UnitID{
						DNSSDName: "wsd",
						UUID:      uu,
						Backend:   mockBackend{},
					},
					Addrs: addr,
				}},
			},
			want: 2,
		},
	}

	for _, test := range tests {
		units := countUnits(test.groups)
		result := out.mergeUnitGrousByUUID(test.groups)

		if len(result) != test.want {
			t.Errorf("%s: %d groups returned, expected %d",
				test.name, len(result), test.want)
		}

		if n := countUnits(result); n != units {
			t.Errorf("%s: %d units returned, expected %d",
				test.name, n, units)
		}
	}
}
