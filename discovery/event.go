// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovery events

package discovery

// Event is the common interface for all events
type Event interface {
	// Name returns the Event name.
	Name() string

	// GetID returns UnitID this Event related to.
	GetID() UnitID
}

var (
	_ = Event(&EventAddUnit{})
	_ = Event(&EventDelUnit{})
	_ = Event(&EventPrinterParameters{})
	_ = Event(&EventScannerParameters{})
	_ = Event(&EventFaxoutParameters{})
	_ = Event(&EventAddEndpoint{})
	_ = Event(&EventDelEndpoint{})
)

// EventAddUnit generated when new print or scan unit is discovered.
//
// Note, if printer has multiple print queues, each queue must be
// announced as a separate unit with the separate ID. Backend may use
// UnitID.SubRealm to make IDs of these unit distinguishable.
type EventAddUnit struct {
	ID UnitID // Unit identity
}

// Name returns the Event name.
func (*EventAddUnit) Name() string {
	return "add-unit"
}

// GetID returns the UnitID this event related to.
func (evnt *EventAddUnit) GetID() UnitID {
	return evnt.ID
}

// EventDelUnit generated when previously announced unit is not
// longer available.
type EventDelUnit struct {
	ID UnitID // Unit identity
}

// Name returns the Event name.
func (*EventDelUnit) Name() string {
	return "del-unit"
}

// GetID returns the UnitID this event related to.
func (evnt *EventDelUnit) GetID() UnitID {
	return evnt.ID
}

// EventPrinterParameters generated when printer parameters
// becomes available or updated.
//
// Backend responsibilities:
//   - Unit MUST exist (i.e., previously announced with the
//     EventAddUnit event and not revoked with the EventDelUnit event)
type EventPrinterParameters struct {
	ID              UnitID            // Unit identity
	TXT             []string          // DND-SD TXT record (if available)
	MakeModel       string            // Manufacturer + Model
	Location        string            // E.g., E.g., "2nd Floor Lab"
	AdminURL        string            // Device administration URL
	IconURL         string            // Device icon URL
	PPDManufacturer string            // I.e., "Hewlett Packard" or "Canon"
	PPDModel        string            // Model name
	Printer         PrinterParameters // Printer parameters
}

// Name returns the Event name.
func (*EventPrinterParameters) Name() string {
	return "printer-parameters"
}

// GetID returns the UnitID this event related to.
func (evnt *EventPrinterParameters) GetID() UnitID {
	return evnt.ID
}

// EventScannerParameters generated when scanner parameters
// become available or updated.
//
// Backend responsibilities:
//   - Unit MUST exist
type EventScannerParameters struct {
	ID        UnitID            // Unit identity
	TXT       []string          // DND-SD TXT record (if available)
	MakeModel string            // Manufacturer + Model
	Location  string            // E.g., E.g., "2nd Floor Lab"
	AdminURL  string            // Device administration URL
	IconURL   string            // Device icon URL
	Scanner   ScannerParameters // Scanner parameters
}

// Name returns the Event name.
func (*EventScannerParameters) Name() string {
	return "scanner-parameters"
}

// GetID returns the UnitID this event related to.
func (evnt *EventScannerParameters) GetID() UnitID {
	return evnt.ID
}

// EventFaxoutParameters generated when faxout parameters
// become available or updated.
//
// Backend responsibilities:
//   - Unit MUST exist
type EventFaxoutParameters struct {
	ID              UnitID            // Unit identity
	TXT             []string          // DND-SD TXT record (if available)
	MakeModel       string            // Manufacturer + Model
	Location        string            // E.g., E.g., "2nd Floor Lab"
	AdminURL        string            // Device administration URL
	IconURL         string            // Device icon URL
	PPDManufacturer string            // I.e., "Hewlett Packard" or "Canon"
	PPDModel        string            // Model name
	Faxout          PrinterParameters // Faxout parameters (the same as printer)
}

// Name returns the Event name.
func (*EventFaxoutParameters) Name() string {
	return "scanner-parameters"
}

// GetID returns the UnitID this event related to.
func (evnt *EventFaxoutParameters) GetID() UnitID {
	return evnt.ID
}

// EventAddEndpoint is generated to report each discovered endpoint.
//
// Backend responsibilities:
//   - Unit MUST exist
//   - The same endpoint MUST NOT be added multiple times.
type EventAddEndpoint struct {
	ID       UnitID // Unit identity
	Endpoint string // URLs of added endpoints
}

// Name returns the Event name.
func (*EventAddEndpoint) Name() string {
	return "add-endpoint"
}

// GetID returns the UnitID this event related to.
func (evnt *EventAddEndpoint) GetID() UnitID {
	return evnt.ID
}

// EventDelEndpoint is generated, when some of the previously reported
// endpoints is not longer available.
//
// Backend responsibilities:
//   - Unit MUST exist
//   - The removed endpoints MUST exist.
type EventDelEndpoint struct {
	ID       UnitID // Unit identity
	Endpoint string // URLs of removed endpoints
}

// Name returns the Event name.
func (*EventDelEndpoint) Name() string {
	return "del-endpoint"
}

// GetID returns the UnitID this event related to.
func (evnt *EventDelEndpoint) GetID() UnitID {
	return evnt.ID
}
