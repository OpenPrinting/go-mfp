// MFP - Miulti-Function Printers and scanners toolkit
// WSD core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// WSD Message

package wsd

import (
	"bytes"
	"fmt"

	"github.com/alexpevzner/mfp/xmldoc"
)

// Msg represents a WSD protocol message.
type Msg struct {
	Header Header // Message header
	Body   Body   // Message body
}

// DecodeMsg decodes [msg] from the wire representation
func DecodeMsg(data []byte) (m Msg, err error) {
	root, err := xmldoc.Decode(NsMap, bytes.NewReader(data))
	if err == nil {
		m, err = msgFromXML(root)
	}
	return
}

// msgFromXML decodes [msg] from the XML tree
func msgFromXML(root xmldoc.Element) (m Msg, err error) {
	const (
		rootName = NsSOAP + ":" + "Envelope"
		hdrName  = NsSOAP + ":" + "Header"
		bodyName = NsSOAP + ":" + "Body"
	)

	defer func() { err = xmlErrWrap(root, err) }()

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
		err = xmlErrMissed(missed.Name)
		return
	}

	// Decode message header
	m.Header, err = DecodeHeader(hdr.Elem)
	if err != nil {
		return
	}

	// Fetch body Element
	name := m.Header.Action.bodyname()
	var elem xmldoc.Element

	if name != "" {
		var ok bool
		elem, ok = body.Elem.ChildByName(name)
		if !ok {
			err = xmlErrWrap(body.Elem, xmlErrMissed(name))
			return
		}
	}

	// Decode message body
	switch m.Header.Action {
	case ActHello:
		m.Body, err = DecodeHello(elem)
	case ActBye:
		m.Body, err = DecodeBye(elem)
	case ActProbe:
		m.Body, err = DecodeProbe(elem)
	case ActProbeMatches:
		m.Body, err = DecodeProbeMatches(elem)
	case ActResolve:
		m.Body, err = DecodeResolve(elem)
	case ActResolveMatches:
		m.Body, err = DecodeResolveMatches(elem)
	case ActGet:
		m.Body, err = DecodeGet(elem)
	case ActGetResponse:
		m.Body, err = DecodeMetadata(elem)
	default:
		err = fmt.Errorf("%s: unhanded action ", m.Header.Action)
		return
	}

	return
}

// Encode encodes [Msg] into its wire representation.
func (m Msg) Encode() []byte {
	buf := bytes.Buffer{}
	m.ToXML().Encode(&buf, NsMap)
	return buf.Bytes()
}

// Format formats [Msg] for logging/
func (m Msg) Format() string {
	return m.ToXML().EncodeIndentString(NsMap, "  ")
}

// ToXML generates XML tree for the message
func (m Msg) ToXML() xmldoc.Element {
	var body []xmldoc.Element
	if bodydata := m.Body.ToXML(); !bodydata.IsZero() {
		body = []xmldoc.Element{bodydata}
	}

	elm := xmldoc.Element{
		Name: NsSOAP + ":" + "Envelope",
		Children: []xmldoc.Element{
			m.Header.ToXML(),
			xmldoc.Element{
				Name:     NsSOAP + ":" + "Body",
				Children: body,
			},
		},
	}

	return elm
}

// MarkUsedNamespace marks [xmldoc.Namespace] entries used by
// data elements within the message body, if any.
//
// This function should not care about Namespace entries, used
// by XML tags: they are handled automatically.
func (m Msg) MarkUsedNamespace(ns xmldoc.Namespace) {
	m.Body.MarkUsedNamespace(ns)
}
