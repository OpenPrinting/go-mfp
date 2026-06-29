// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// WS-Scan SOAP message

package wsscan

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/OpenPrinting/go-mfp/log/trace"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/uuid"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Message represents a WS-Scan SOAP message, consisting of
// a [Header] and a [Body].
type Message struct {
	Header Header
	Body   Body
}

// DecodeMessage decodes a [Message] from the XML tree.
func DecodeMessage(root xmldoc.Element) (msg Message, err error) {
	return decodeMessageXML(root)
}

// decodeMessageXML decodes a [Message] from the XML tree.
func decodeMessageXML(root xmldoc.Element) (msg Message, err error) {
	const (
		rootName = NsSOAP + ":" + "Envelope"
		hdrName  = NsSOAP + ":" + "Header"
		bodyName = NsSOAP + ":" + "Body"
	)

	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	// Check root element
	if root.Name != rootName {
		err = fmt.Errorf("%s: missed", rootName)
		return
	}

	// Look for Header and Body elements
	hdr := xmldoc.Lookup{Name: hdrName, Required: true}
	body := xmldoc.Lookup{Name: bodyName, Required: true}

	missed := root.Lookup(&hdr, &body)
	if missed != nil {
		err = xmldoc.XMLErrMissed(missed.Name)
		return
	}

	// Decode message header
	msg.Header, err = decodeHeader(hdr.Elem)
	if err != nil {
		return
	}

	// Find the body child element by action
	bodyChildName := msg.Header.Action.bodyElementName()
	if bodyChildName == "" {
		err = fmt.Errorf("unknown action: %s", msg.Header.Action)
		return
	}

	child, ok := body.Elem.ChildByName(bodyChildName)
	if !ok {
		err = xmldoc.XMLErrMissed(bodyChildName)
		err = xmldoc.XMLErrWrap(body.Elem, err)
		return
	}

	// Decode message body by action
	switch msg.Header.Action {
	case ActGetScannerElements:
		v, e := decodeGetScannerElementsRequest(child)
		msg.Body, err = &v, e
	case ActGetScannerElementsResponse:
		v, e := decodeGetScannerElementsResponse(child)
		msg.Body, err = &v, e
	case ActCreateScanJob:
		v, e := decodeCreateScanJobRequest(child)
		msg.Body, err = &v, e
	case ActCreateScanJobResponse:
		v, e := decodeCreateScanJobResponse(child)
		msg.Body, err = &v, e
	case ActRetrieveImage:
		v, e := decodeRetrieveImageRequest(child)
		msg.Body, err = &v, e
	case ActRetrieveImageResponse:
		v, e := decodeRetrieveImageResponse(child)
		msg.Body, err = &v, e
	case ActCancelJob:
		v, e := decodeCancelJobRequest(child)
		msg.Body, err = &v, e
	case ActCancelJobResponse:
		msg.Body = &CancelJobResponse{}
	case ActGetActiveJobs:
		msg.Body = &GetActiveJobsRequest{}
	case ActGetActiveJobsResponse:
		v, e := decodeGetActiveJobsResponse(child)
		msg.Body, err = &v, e
	case ActGetJobHistory:
		msg.Body = &GetJobHistoryRequest{}
	case ActGetJobHistoryResponse:
		v, e := decodeGetJobHistoryResponse(child)
		msg.Body, err = &v, e
	default:
		err = fmt.Errorf("unhandled action: %s", msg.Header.Action)
	}

	return
}

// Encode encodes the [Message] into its wire representation.
func (msg Message) Encode() []byte {
	buf := bytes.Buffer{}
	ns := generic.CopySlice(NsMap)
	msg.toXML().Encode(&buf, ns)
	return buf.Bytes()
}

// Format formats the [Message] for logging.
func (msg Message) Format() string {
	ns := generic.CopySlice(NsMap)
	return msg.toXML().EncodeIndentString(ns, "  ")
}

// toXML generates the XML tree for the SOAP envelope.
func (msg Message) toXML() xmldoc.Element {
	var bodyChildren []xmldoc.Element
	if msg.Body != nil {
		bodyChildren = []xmldoc.Element{msg.Body.ToXML()}
	}

	return xmldoc.Element{
		Name: NsSOAP + ":" + "Envelope",
		Children: []xmldoc.Element{
			msg.Header.toXML(),
			{
				Name:     NsSOAP + ":" + "Body",
				Children: bodyChildren,
			},
		},
	}
}

// mtomContentType returns the Content-Type header value for the
// MTOM/XOP multipart encoding using the given boundary and
// envelope Content-ID.
func mtomContentType(boundary, envelopeCID string) string {
	return fmt.Sprintf(
		`multipart/related;`+
			` type="application/xop+xml";`+
			` boundary="%s";`+
			` start="<%s>";`+
			` start-info="application/soap+xml"`,
		boundary, envelopeCID)
}

// writeMTOM encodes the message as an MTOM/XOP multipart response,
// with the SOAP envelope as the first part and the binary
// attachment from [RetrieveImageResponse] as the second part.
//
// Message.Body type must be *RetrieveImageResponse
func (msg Message) writeMTOM(query *transport.ServerQuery) error {
	body := msg.Body.(*RetrieveImageResponse)

	// Write HTTP header
	boundary := uuid.Random().String()
	envelopeCID := uuid.Random().String()

	query.ResponseHeader().Set("Content-Type",
		mtomContentType(boundary, envelopeCID))
	query.WriteHeader(http.StatusOK)

	// Encode the SOAP envelope
	soapData := msg.Encode()

	// Create multipart writer with the pre-determined boundary
	mw := multipart.NewWriter(query)
	if err := mw.SetBoundary(boundary); err != nil {
		return err
	}

	// Part 1: SOAP envelope
	soapPart, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type": {`application/xop+xml;` +
			` charset=UTF-8;` +
			` type="application/soap+xml"`},
		"Content-Transfer-Encoding": {"binary"},
		"Content-Id":                {"<" + envelopeCID + ">"},
	})
	if err != nil {
		return err
	}

	if _, err = soapPart.Write(soapData); err != nil {
		return err
	}

	// Notify tracer. We do it now, so only the attached
	// image will be captured by tracer as data attachment.
	image := trace.OnResponse(query, traceMessage{msg}, body.Image)
	defer image.Close()

	// Part 2: Image data (streamed)
	imagePart, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {body.ContentType},
		"Content-Transfer-Encoding": {"binary"},
		"Content-Id":                {"<" + body.ScanData.ContentID + ">"},
	})
	if err != nil {
		return err
	}
	if _, err = io.Copy(imagePart, image); err != nil {
		return err
	}

	return mw.Close()
}
