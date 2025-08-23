// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// condition history entry

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ConditionHistoryEntry represents the <wscn:ConditionHistoryEntry> element,
// providing details about a previously active condition that has been cleared.
type ConditionHistoryEntry struct {
	ClearTime DateTime
	Component Component
	Name      NameElement
	Severity  Severity
	Time      DateTime
}

// toXML generates XML tree for the [ConditionHistoryEntry].
func (che ConditionHistoryEntry) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		che.ClearTime.toXML(NsWSCN + ":ClearTime"),
		che.Component.toXML(NsWSCN + ":Component"),
		che.Name.toXML(NsWSCN + ":Name"),
		che.Severity.toXML(NsWSCN + ":Severity"),
		che.Time.toXML(NsWSCN + ":Time"),
	}
	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeConditionHistoryEntry decodes [ConditionHistoryEntry] from the XML tree.
func decodeConditionHistoryEntry(root xmldoc.Element) (
	che ConditionHistoryEntry, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	clearTime := xmldoc.Lookup{
		Name:     NsWSCN + ":ClearTime",
		Required: true,
	}
	component := xmldoc.Lookup{
		Name:     NsWSCN + ":Component",
		Required: true,
	}
	name := xmldoc.Lookup{
		Name:     NsWSCN + ":Name",
		Required: true,
	}
	severity := xmldoc.Lookup{
		Name:     NsWSCN + ":Severity",
		Required: true,
	}
	time := xmldoc.Lookup{
		Name:     NsWSCN + ":Time",
		Required: true,
	}

	missed := root.Lookup(
		&clearTime,
		&component,
		&name,
		&severity,
		&time,
	)
	if missed != nil {
		return che, xmldoc.XMLErrMissed(missed.Name)
	}

	if che.ClearTime, err = decodeDateTime(clearTime.Elem); err != nil {
		return che, fmt.Errorf("clearTime: %w", err)
	}
	if che.Component, err = decodeComponent(component.Elem); err != nil {
		return che, fmt.Errorf("component: %w", err)
	}
	if che.Name, err = decodeNameElement(name.Elem); err != nil {
		return che, fmt.Errorf("name: %w", err)
	}
	if che.Severity, err = decodeSeverity(severity.Elem); err != nil {
		return che, fmt.Errorf("severity: %w", err)
	}
	if che.Time, err = decodeDateTime(time.Elem); err != nil {
		return che, fmt.Errorf("time: %w", err)
	}

	return che, nil
}
