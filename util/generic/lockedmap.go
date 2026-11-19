// MFP - Miulti-Function Printers and scanners toolkit
// Useful generics
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Goroutine-safe map

package generic

import "sync"

// LockedMap adds goroutine-safety to the standard go map.
type LockedMap[K comparable, V any] struct {
	data map[K]V    // Underlying map
	lock sync.Mutex // Access lock
}

// NewLockedMap creates a new LockedMap
func NewLockedMap[K comparable, V any]() *LockedMap[K, V] {
	return &LockedMap[K, V]{
		data: make(map[K]V),
	}
}

// Clear clears the map.
func (m *LockedMap[K, V]) Clear() {
	m.lock.Lock()
	clear(m.data)
	m.lock.Unlock()
}

// Count returns count of elements in the map.
func (m *LockedMap[K, V]) Count() (count int) {
	m.lock.Lock()
	count = len(m.data)
	m.lock.Unlock()
	return count
}

// Delete deletes the element with the specified key.
func (m *LockedMap[K, V]) Delete(key K) {
	m.lock.Lock()
	delete(m.data, key)
	m.lock.Unlock()
}

// Empty reports if map is empty.
func (m *LockedMap[K, V]) Empty() bool {
	return m.Count() == 0
}

// Load returns the element with the specified key or the
// zero value of the element's type and indication that
// element was found in the map.
func (m *LockedMap[K, V]) Load(key K) (val V, ok bool) {
	m.lock.Lock()
	val, ok = m.data[key]
	m.lock.Unlock()

	return
}

// Store stores element under the specified key.
func (m *LockedMap[K, V]) Store(key K, val V) {
	m.lock.Lock()
	m.data[key] = val
	m.lock.Unlock()
}
