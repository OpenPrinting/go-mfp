// MFP - Miulti-Function Printers and scanners toolkit
// WSD core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// AppSequence

package wsd

import (
	"errors"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// AppSequence provides a mechanism that allows a receiver
// to order messages that may have been received out of order.
//
// It is included into the announcement and response messages
// ([Hello], [Bye], [ProbeMatches], and [ResolveMatches]).
type AppSequence struct {
	InstanceID    uint64               // MUST increment on each reboot
	MessageNumber uint64               // MUST increment on each message
	SequenceID    optional.Val[AnyURI] // Sequence within instance
}

// DecodeAppSequence decodes AppSequence from the XML tree
func DecodeAppSequence(root xmldoc.Element) (seq AppSequence, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	instanceID := xmldoc.LookupAttr{Name: "InstanceId", Required: true}
	messageNumber := xmldoc.LookupAttr{Name: "MessageNumber", Required: true}
	sequenceID := xmldoc.LookupAttr{Name: "SequenceId"}

	missed := root.LookupAttrs(&instanceID, &messageNumber, &sequenceID)
	if missed != nil {
		err = errors.New("missed attribyte")
		err = xmldoc.XMLErrWrapName("@"+missed.Name, err)
		err = xmldoc.XMLErrWrap(root, err)
		return
	}

	seq.InstanceID, err = decodeUint64Attr(instanceID.Attr)
	if err == nil {
		seq.MessageNumber, err = decodeUint64Attr(messageNumber.Attr)
	}

	if err == nil && sequenceID.Found {
		var id AnyURI
		id, err = DecodeAnyURIAttr(sequenceID.Attr)
		if err == nil {
			seq.SequenceID = optional.New(id)
		}
	}

	return
}

// ToXML generates XML tree for the AppSequence
func (seq AppSequence) ToXML() xmldoc.Element {
	elm := xmldoc.Element{
		Name: NsDiscovery + ":" + "AppSequence",
		Attrs: []xmldoc.Attr{
			{
				Name:  "InstanceId",
				Value: strconv.FormatUint(seq.InstanceID, 10),
			},
			{
				Name:  "MessageNumber",
				Value: strconv.FormatUint(seq.MessageNumber, 10),
			},
		},
	}

	if seq.SequenceID != nil {
		elm.Attrs = append(elm.Attrs, xmldoc.Attr{
			Name:  "SequenceId",
			Value: string(*seq.SequenceID),
		})
	}

	return elm
}
