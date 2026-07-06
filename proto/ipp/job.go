// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Job state

package ipp

import (
	"strings"
	"sync"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// job represents state of the job
type job struct {
	JobStatus                     // Job status attributes
	JobCreateOperation            // Job create-time operation attributes
	JobAttributes                 // Job creation attributes
	SendDocumentActive bool       // Send-Document in progress
	cancelPending      bool       // Cancel-Job accepted, not yet canceled
	lock               sync.Mutex // Access lock
}

// newJob creates a new job.
func newJob(ops *JobCreateOperation, attrs *JobAttributes) *job {
	uu := uuid.Random()
	uri := strings.Join([]string{ops.PrinterURI, "jobs", uu.String()}, "/")

	j := &job{
		JobStatus: JobStatus{
			JobImpressionsCompleted: optional.New(0),
			JobMediaSheetsCompleted: optional.New(0),
			JobName:                 ops.JobName,
			JobOriginatingUserName:  ops.RequestingUserName,
			JobState:                EnJobStatePendingHeld,
			JobStateReasons:         []KwJobStateReasons{KwJobStateReasonsJobIncoming},
			JobURI:                  uri,
		},
		JobCreateOperation: *ops,
		JobAttributes:      *attrs,
	}

	return j
}

// Lock acquires the job's mutex
func (j *job) Lock() {
	j.lock.Lock()
}

// Unlock releases the job's mutex
func (j *job) Unlock() {
	j.lock.Unlock()
}
