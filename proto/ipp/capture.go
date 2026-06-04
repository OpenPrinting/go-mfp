// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// IPP document capture hook.

package ipp

// DocumentReceiver is called when a document is received via Send-Document.
//
// jobID is the IPP job ID assigned by the printer.
// format is the MIME type of the document (e.g., "application/pdf",
// "image/pwg-raster").
// data contains the complete document bytes.
type DocumentReceiver func(jobID int, format string, data []byte)

// SetDocumentReceiver installs fn as the document receiver callback.
// fn is called once per Send-Document request after the body is fully read.
// Pass nil to clear a previously installed receiver.
func (printer *Printer) SetDocumentReceiver(fn DocumentReceiver) {
	printer.receiver = fn
}
