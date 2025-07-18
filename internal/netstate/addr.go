// MFP - Miulti-Function Printers and scanners toolkit
// Network state monitoring
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Network Interface addresses

package netstate

import (
	"net"
	"net/netip"
)

// Addr represents a single IP address with a mask, assigned to a network
// interface.
//
// Unlike [net.IP] or [net.IPAddr], Addr is a comparable value type
// (it supports == and can be used as a map key) and is immutable.
//
// An interface may have multiple addresses that can belong to the same
// or different IP networks. Addresses are grouped by their IP network
// membership. Within each group, exactly one address is marked as Primary.
//
// In other words:
//   - If all interface addresses belong to different IP networks, they will
//     all be marked as Primary.
//   - If multiple interface addresses belong to the same IP network, only one
//     of them will be chosen as Primary.
//
// Two addresses are considered to belong to the same IP network if their
// address ranges (taking the mask into account) overlap. You can use
// [Addr.Overlaps] to test whether any two addresses overlap.
//
// Strictly speaking, ranges covered by two overlapping addresses either
// equal, if masks are the same, or nest, if mask of the "inner" address
// is narrower that mask of the "outer" address.
//
// Use [Addr.Narrower] to determine which of two overlapping addresses
// has the narrower mas
type Addr struct {
	netip.Prefix       // IP address with mask
	nif          NetIf // Interface that owns the address
}

// AddrFromIPNet makes address from the [net.IPNet]
func AddrFromIPNet(ipn net.IPNet, nif NetIf) Addr {
	ip, _ := netip.AddrFromSlice(ipn.IP)
	ip = ip.Unmap()
	bits, _ := ipn.Mask.Size()
	prefix := netip.PrefixFrom(ip, bits)
	return Addr{prefix, nif}
}

// Addr returns IP address.
//
// If address is a link-local IPv6 address, it comes with the
// appropriate zone.
func (addr Addr) Addr() netip.Addr {
	ip := addr.Prefix.Addr()
	if ip.Is6() && ip.IsLinkLocalUnicast() {
		ip = ip.WithZone(addr.nif.Name())
	}
	return ip
}

// Interface returns the network interface that owns the address.
func (addr Addr) Interface() NetIf {
	return addr.nif
}

// SameInterface reports if two addresses belong to the same
// network interface.
//
// Note, we consider two interfaces equal if they have equal
// [net.Interface.Index] and [net.Interface.Name]. Other parts
// of the [net.Interface] considered interface parameters, not
// interface identity.
func (addr Addr) SameInterface(addr2 Addr) bool {
	return addr.Interface() == addr2.Interface()
}

// Less orders [Addr] for sorting.
//
// The sorting order is following:
//
//   - if addresses belongs to different interfaces, they are
//     sorted by [net.Interface.Index], in acceding order
//   - if interface indices are the same, but name differ, addresses
//     are sorted by interface name, in acceding order
//   - otherwise, if addresses belong to the different address
//     families, they are sorted by address family, IPv4 first
//   - otherwise, if IP addresses not the same, they are sorted
//     by IP address, in lexicographical acceding order
//   - otherwise, if masks are different, addresses are sorted by
//     network mask, the narrowed first
//   - otherwise, addresses are equal
func (addr Addr) Less(addr2 Addr) bool {
	switch {
	case !addr.SameInterface(addr2):
		// Sort by net.Interface.Index
		return addr.Interface().Less(addr2.Interface())
	case addr.Is4() != addr2.Is4():
		// Sort by address family, IP4 first
		return addr.Is4()
	case addr.Addr() != addr2.Addr():
		// Sort by IP address, lexicographically
		return addr.Addr().Less(addr2.Addr())
	default:
		// Sort by network mask, the narrowed first.
		return addr.Narrower(addr2)
	}
}

// Unmasked returns the [Addr] with the same IP address but mask
// that corresponds to single IP:
//
//	127.0.0.1/24 -> 127.0.0.1/32
//	"::1/24" -> ""::1/128"
func (addr Addr) Unmasked() Addr {
	ip := addr.Addr()
	bits := ip.BitLen()
	prefix := netip.PrefixFrom(ip, bits)
	return Addr{prefix, addr.nif}
}

// Similar reports whether two addresses are the same, ignoring
// difference in address mask
func (addr Addr) Similar(addr2 Addr) bool {
	return addr.Unmasked() == addr2.Unmasked()
}

// Overlaps reports whether two addresses overlap.
//
// Overlapped addressed are addresses for which all following is true:
//   - they belong to the same network interface
//   - they belong to the same address family, both either IP4 or IP6
//   - their address range, taking Mask into account, overlap
func (addr Addr) Overlaps(addr2 Addr) bool {
	var answer bool
	if addr.SameInterface(addr2) {
		answer = addr.Prefix.Overlaps(addr2.Prefix)
	}
	return answer
}

// Narrower reports whether addr is narrower that addr2.
//
// It means the following:
//   - addr and addr2 overlap (see [Addr.Overlap] for definition).
//   - mask of addr is narrower (contains more leading ones and less
//     trailing zeroes) that mask of addr2
func (addr Addr) Narrower(addr2 Addr) bool {
	var answer bool
	if addr.Overlaps(addr2) {
		answer = addr.Bits() > addr2.Bits()
	}
	return answer
}

// Wider is the opposite to the [Addr.Narrower]
func (addr Addr) Wider(addr2 Addr) bool {
	return addr2.Narrower(addr)
}

// Is4 tells is [Addr] is IP4 address.
func (addr Addr) Is4() bool {
	return addr.Prefix.Addr().Is4() || addr.Prefix.Addr().Is4In6()
}

func (addr Addr) String() string {
	prefix := netip.PrefixFrom(addr.Prefix.Addr().Unmap(),
		addr.Prefix.Bits())
	return prefix.String()
}
