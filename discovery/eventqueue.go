// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Event queue

package discovery

import (
	"context"
	"sync"
)

// Eventqueue represents a queue of [Event].
//
// [Backend] communicates with discovery system by pushing events
// into the queue. The queue is created and owned by the discovery
// system and passed as parameter to the Backend constructor when
// Backend is being created.
//
// See description of each particular Event for
// Backend's responsibility when generating this kind of Event.
type Eventqueue struct {
	events    []Event       // Events in the queue
	readychan chan struct{} // Signaled when more events is available
	lock      sync.Mutex    // Access lock
}

// NewEventqueue creates the new Eventqueue
func NewEventqueue() *Eventqueue {
	return &Eventqueue{
		events:    make([]Event, 0, 32),
		readychan: make(chan struct{}, 1),
	}
}

// Push pushes event into the queue.
func (q *Eventqueue) Push(e Event) {
	q.lock.Lock()
	q.events = append(q.events, e)
	q.lock.Unlock()

	select {
	case q.readychan <- struct{}{}:
	default:
	}
}

// Count returns the current queue length.
func (q *Eventqueue) Count() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.events)
}

// pull returns next event out of the queue.
//
// If queue is empty, it will wait until more events is available
// or Context is expired.
//
// The only case when error is returned is caused by the
// Context expiration.
func (q *Eventqueue) pull(ctx context.Context) (Event, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for ctx.Err() == nil {
		// Peek next event, if available
		if len(q.events) != 0 {
			e := q.events[0]
			copy(q.events, q.events[1:])
			q.events = q.events[:len(q.events)-1]
			return e, nil
		}

		// Wait for the more events or context expiration
		q.lock.Unlock()
		select {
		case <-q.readychan:
		case <-ctx.Done():
		}
		q.lock.Lock()
	}

	return nil, ctx.Err()
}
