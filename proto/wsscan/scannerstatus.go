// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scanner status

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ScannerStatus represents the <wscn:ScannerStatus> element,
// providing comprehensive information about the scanner's current state.
type ScannerStatus struct {
	ActiveConditions    []DeviceCondition
	ConditionHistory    []ConditionHistoryEntry
	ScannerCurrentTime  DateTime
	ScannerState        ScannerState
	ScannerStateReasons []ScannerStateReason
}

// toXML generates XML tree for the [ScannerStatus].
func (ss ScannerStatus) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{}

	// ActiveConditions slice
	if len(ss.ActiveConditions) > 0 {
		acChildren := make([]xmldoc.Element, len(ss.ActiveConditions))
		for i, v := range ss.ActiveConditions {
			acChildren[i] = v.toXML(NsWSCN + ":DeviceCondition")
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":ActiveConditions",
			Children: acChildren,
		})
	}

	// ConditionHistory (optional)
	if len(ss.ConditionHistory) > 0 {
		chChildren := make([]xmldoc.Element, len(ss.ConditionHistory))
		for i, v := range ss.ConditionHistory {
			chChildren[i] = v.toXML(NsWSCN + ":ConditionHistoryEntry")
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":ConditionHistory",
			Children: chChildren,
		})
	}

	// ScannerCurrentTime
	children = append(children, ss.ScannerCurrentTime.toXML(
		NsWSCN+":ScannerCurrentTime"))

	// ScannerState
	children = append(children, ss.ScannerState.toXML(NsWSCN+":ScannerState"))

	// ScannerStateReasons slice
	if len(ss.ScannerStateReasons) > 0 {
		ssrChildren := make([]xmldoc.Element, len(ss.ScannerStateReasons))
		for i, v := range ss.ScannerStateReasons {
			ssrChildren[i] = v.toXML(NsWSCN + ":ScannerStateReason")
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":ScannerStateReasons",
			Children: ssrChildren,
		})
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeScannerStatus decodes [ScannerStatus] from the XML tree.
func decodeScannerStatus(root xmldoc.Element) (ss ScannerStatus, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	// Lookup all required XML elements
	scannerCurrentTime := xmldoc.Lookup{
		Name:     NsWSCN + ":ScannerCurrentTime",
		Required: true,
	}
	scannerState := xmldoc.Lookup{
		Name:     NsWSCN + ":ScannerState",
		Required: true,
	}

	// Lookup optional XML elements
	activeConditions := xmldoc.Lookup{
		Name:     NsWSCN + ":ActiveConditions",
		Required: true,
	}
	conditionHistory := xmldoc.Lookup{
		Name:     NsWSCN + ":ConditionHistory",
		Required: false,
	}
	scannerStateReasons := xmldoc.Lookup{
		Name:     NsWSCN + ":ScannerStateReasons",
		Required: true,
	}

	missed := root.Lookup(
		&scannerCurrentTime,
		&scannerState,
		&activeConditions,
		&conditionHistory,
		&scannerStateReasons,
	)
	if missed != nil {
		return ss, xmldoc.XMLErrMissed(missed.Name)
	}

	// Required fields
	if ss.ScannerCurrentTime, err = decodeDateTime(scannerCurrentTime.Elem); err != nil {
		return ss, fmt.Errorf("scannerCurrentTime: %w", err)
	}
	if ss.ScannerState, err = decodeScannerState(scannerState.Elem); err != nil {
		return ss, fmt.Errorf("scannerState: %w", err)
	}

	// ActiveConditions slice
	if activeConditions.Elem.Children != nil {
		for _, child := range activeConditions.Elem.Children {
			if child.Name == NsWSCN+":DeviceCondition" {
				val, err := decodeDeviceCondition(child)
				if err != nil {
					return ss, fmt.Errorf("ActiveConditions: "+
						"invalid DeviceCondition: %w", err)
				}
				ss.ActiveConditions = append(ss.ActiveConditions, val)
			}
		}
	}

	// Optional ConditionHistory slice
	if conditionHistory.Elem.Children != nil {
		var ch []ConditionHistoryEntry
		for _, child := range conditionHistory.Elem.Children {
			if child.Name == NsWSCN+":ConditionHistoryEntry" {
				val, err := decodeConditionHistoryEntry(child)
				if err != nil {
					return ss, fmt.Errorf("ConditionHistory: "+
						"invalid ConditionHistoryEntry: %w", err)
				}
				ch = append(ch, val)
			}
		}
		if len(ch) > 0 {
			ss.ConditionHistory = ch
		}
	}

	// ScannerStateReasons slice
	if scannerStateReasons.Elem.Children != nil {
		for _, child := range scannerStateReasons.Elem.Children {
			if child.Name == NsWSCN+":ScannerStateReason" {
				val, err := decodeScannerStateReason(child)
				if err != nil {
					return ss, fmt.Errorf("ScannerStateReasons: "+
						"invalid ScannerStateReason: %w", err)
				}
				ss.ScannerStateReasons = append(ss.ScannerStateReasons, val)
			}
		}
	}

	return ss, nil
}
