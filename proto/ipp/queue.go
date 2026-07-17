// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Queue of jobs

package ipp

import (
	"math"
	"sync"
)

// queue manages queue of jobs
type queue struct {
	lock   sync.Mutex      // Access lock
	nextid int32           // Next JOB ID
	jobs   []*job          // Queued jobs
	byID   map[int]*job    // Job by JobID
	byURI  map[string]*job // Job by JobURI
}

// newQueue creates new queue.
func newQueue() *queue {
	q := &queue{
		nextid: 1,
		byID:   make(map[int]*job),
		byURI:  make(map[string]*job),
	}
	return q
}

// Push pushes new job into the queue.
// It assigns the job.JobId
func (q *queue) Push(j *job) {
	q.lock.Lock()
	defer q.lock.Unlock()

	j.JobID = q.allocJobID()

	q.jobs = append(q.jobs, j)
	q.byID[j.JobID] = j
	q.byURI[j.JobURI] = j
}

// JobByID returns job by its ID
func (q *queue) JobByID(id int) *job {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.byID[id]
}

// JobByURI returns job by its URI
func (q *queue) JobByURI(uri string) *job {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.byURI[uri]
}

// Jobs returns a snapshot of queued jobs in submission order.
func (q *queue) Jobs() []*job {
	q.lock.Lock()
	defer q.lock.Unlock()

	out := make([]*job, len(q.jobs))
	copy(out, q.jobs)
	return out
}

// allocJobID allocates the next JobID.
// It must be called under q.lock.
func (q *queue) allocJobID() int {
	for {
		id := int(q.nextid)
		if q.nextid == math.MaxInt32 {
			q.nextid = 1
		} else {
			q.nextid++
		}

		if q.byID[id] == nil {
			return id
		}
	}
}
