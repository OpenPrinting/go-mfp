// MFP - Miulti-Function Printers and scanners toolkit
// Useful generics
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Generic sets

package generic

// Set is the generic set of any comparable objects.
//
// Set cannot be simultaneously accessed from multiple goroutines.
// If you need goroutine safety, use [LockedSet].
type Set[T comparable] struct {
	members map[T]struct{} // Members of the set
}

// NewSet creates a new Set
func NewSet[T comparable]() Set[T] {
	return Set[T]{
		members: make(map[T]struct{}),
	}
}

// Clear purges the set
func (s Set[T]) Clear() {
	for member := range s.members {
		delete(s.members, member)
	}
}

// Empty reports if set is empty
func (s Set[T]) Empty() bool {
	return len(s.members) == 0
}

// Contains reports if member already in the set
func (s Set[T]) Contains(member T) bool {
	_, found := s.members[member]
	return found
}

// Add adds member to the set
func (s Set[T]) Add(member T) {
	s.members[member] = struct{}{}
}

// Del deletes member from the set
func (s Set[T]) Del(member T) {
	delete(s.members, member)
}

// TestAndAdd adds member to the set and returns true if it was actually added.
func (s Set[T]) TestAndAdd(member T) (added bool) {
	if !s.Contains(member) {
		s.Add(member)
		added = true
	}
	return
}

// TestAndDel deletes member from the set and returns true if it was actually
// deleted.
func (s Set[T]) TestAndDel(member T) (deleted bool) {
	if s.Contains(member) {
		s.Del(member)
		deleted = true
	}
	return
}

// ForEach applies function to the each member of the set
func (s Set[T]) ForEach(f func(T)) {
	for member := range s.members {
		f(member)
	}
}
