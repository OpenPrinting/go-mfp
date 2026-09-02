// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Avahi glue

package dnssd

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/OpenPrinting/go-avahi"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/internal/zone"
	"github.com/OpenPrinting/go-mfp/util/generic"
)

// Parameters:
const (
	// Avahi client may fail to start if avahi-daemon is not
	// running or for similar reasons.
	//
	// If failed, Avahi client will automatically restart
	// with the following interval between attempts:
	avahiClientRestartInterval = 1 * time.Second
)

// avahiClient wraps avahi.Client and adds some additional
// functionality on a centralized manner.
//
// In particular, it manages tables of avahiService and avahiHostname
// structures that belong to the client and survives over Avahi
// client restart.
type avahiClient struct {
	avahiClnt   *avahi.Client                       // The avahi.Client
	poller      *avahi.Poller                       // The avahi.Poller
	services    map[avahiServiceKey]*avahiService   // Table of services
	hostnames   map[avahiHostnameKey]*avahiHostname // Table of hostnames
	lookupFlags LookupFlags                         // Lookup flags
	domain      string                              // Lookup domain
}

// newAvahiClient creates a new avahiClient
func newAvahiClient(domain string, flags LookupFlags) (
	*avahiClient, error) {

	// Create avahi.Client
	avahiClnt, err := avahi.NewClient(avahi.ClientLoopbackWorkarounds)
	if err != nil {
		return nil, err
	}

	// Create avahiClient structure
	clnt := &avahiClient{
		avahiClnt:   avahiClnt,
		poller:      avahi.NewPoller(),
		services:    make(map[avahiServiceKey]*avahiService),
		hostnames:   make(map[avahiHostnameKey]*avahiHostname),
		lookupFlags: flags,
		domain:      domain,
	}

	clnt.poller.AddClient(avahiClnt)

	return clnt, nil
}

// close closes the client and releases all resources it holds,
func (clnt *avahiClient) Close() {
	clnt.purgeTables()
	clnt.avahiClnt.Close()
}

// restart attempts to restart avahi.Client
func (clnt *avahiClient) Restart(ctx context.Context) error {
	for {
		// Close Avahi client and purge all tables
		clnt.avahiClnt.Close()
		clnt.purgeTables()

		// Pause for avahiClientRestartInterval.
		// Return immediately, if ctx has expired.
		timer := time.NewTimer(avahiClientRestartInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}

		// Try to restart the Client
		flags := avahi.ClientLoopbackWorkarounds
		avahiClnt, err := avahi.NewClient(flags)
		if err == nil {
			clnt.avahiClnt = avahiClnt
			clnt.poller.AddClient(avahiClnt)
			return err
		}
	}
}

// purgeTables purges tables of avahiService and avahiHostname
func (clnt *avahiClient) purgeTables() {
	for _, service := range clnt.services {
		service.Delete()
	}

	if len(clnt.hostnames) != 0 {
		panic("internal error")
	}
}

// poll returns a new event from any of event sources owned
// by the avahiClient.
func (clnt *avahiClient) poll(ctx context.Context) (any, error) {
	return clnt.poller.Poll(ctx)
}

// newServiceBrowser creates a new avahi.ServiceBrowser for the
// specified service type.
func (clnt *avahiClient) NewServiceBrowser(svctype string) (
	*avahi.ServiceBrowser, error) {

	browser, err := avahi.NewServiceBrowser(
		clnt.avahiClnt,
		avahi.IfIndexUnspec,
		avahi.ProtocolUnspec,
		svctype,
		clnt.domain,
		avahiLookupFlags(clnt.lookupFlags),
	)

	if err != nil {
		return nil, err
	}

	clnt.poller.AddServiceBrowser(browser)
	return browser, nil
}

// newServiceResolver creates a new avahi.ServiceResolver for
// the parameters specified by the avahiServiceKey.
func (clnt *avahiClient) NewServiceResolver(key avahiServiceKey) (
	*avahi.ServiceResolver, error) {

	flags := avahiLookupFlags(clnt.lookupFlags)
	flags |= avahi.LookupNoAddress | avahi.LookupNoTXT

	resolver, err := avahi.NewServiceResolver(
		clnt.avahiClnt,
		key.IfIdx,
		key.Proto,
		key.InstanceName,
		key.SvcType,
		key.Domain,
		avahi.ProtocolUnspec,
		flags,
	)

	if err != nil {
		return nil, err
	}

	clnt.poller.AddServiceResolver(resolver)
	return resolver, nil
}

// newTxtBrowser creates a new avahi.RecordBrowser for browsing for
// the TXT records with parameters specified by the avahiServiceKey.
func (clnt *avahiClient) NewTxtBrowser(key avahiServiceKey) (
	*avahi.RecordBrowser, error) {

	browser, err := avahi.NewRecordBrowser(
		clnt.avahiClnt,
		key.IfIdx,
		key.Proto,
		key.FQDN(),
		avahi.DNSClassIN,
		avahi.DNSTypeTXT,
		avahiLookupFlags(clnt.lookupFlags),
	)

	if err != nil {
		return nil, err
	}

	clnt.poller.AddRecordBrowser(browser)
	return browser, nil
}

// NewAddrBrowser creates a new avahi.RecordBrowser for browsing for
// the A/AAAA records with parameters specified by the avahiHostnameKey.
func (clnt *avahiClient) NewAddrBrowser(key avahiHostnameKey) (
	*avahi.RecordBrowser, error) {

	rtype := avahi.DNSTypeA
	if key.Proto == avahi.ProtocolIP6 {
		rtype = avahi.DNSTypeAAAA
	}

	browser, err := avahi.NewRecordBrowser(
		clnt.avahiClnt,
		key.IfIdx,
		key.Proto,
		key.Hostname,
		avahi.DNSClassIN,
		rtype,
		avahiLookupFlags(clnt.lookupFlags),
	)

	if err != nil {
		return nil, err
	}

	clnt.poller.AddRecordBrowser(browser)
	return browser, nil
}

// AddService adds a new avahiService.
func (clnt *avahiClient) AddService(key avahiServiceKey,
	svcResolver *avahi.ServiceResolver, txtBrowser *avahi.RecordBrowser) {

	service := &avahiService{
		clnt:        clnt,
		key:         key,
		svcResolver: svcResolver,
		txtBrowser:  txtBrowser,
		units:       make(map[string]*unit),
	}

	clnt.services[key] = service

}

// HasService reports if avahiService already exist for the key
func (clnt *avahiClient) HasService(key avahiServiceKey) bool {
	return clnt.GetService(key) != nil
}

// GetService returns existent avahiService
func (clnt *avahiClient) GetService(key avahiServiceKey) *avahiService {
	return clnt.services[key]
}

// AddService adds a new avahiHostname.
func (clnt *avahiClient) AddHostname(key avahiHostnameKey,
	addrBrowser *avahi.RecordBrowser) *avahiHostname {

	hostname := &avahiHostname{
		clnt:        clnt,
		key:         key,
		addrBrowser: addrBrowser,
		addrs:       generic.NewSet[netip.Addr](),
		services:    generic.NewSet[*avahiService](),
	}

	clnt.hostnames[key] = hostname

	return hostname
}

// GetHostname returns existent avahiHostname
func (clnt *avahiClient) GetHostname(key avahiHostnameKey) *avahiHostname {
	return clnt.hostnames[key]
}

// avahiService is the per-service-instance structure
//
// It is created when new service instance is discovered and
// manages resolving of resource records associated with the
// service: service hostname, TXT record and IP addresses.
type avahiService struct {
	clnt        *avahiClient           // The owner
	key         avahiServiceKey        // Identity
	svcResolver *avahi.ServiceResolver // Service resolver
	txtBrowser  *avahi.RecordBrowser   // TXT record resolver
	port        uint16                 // IP port
	hostname    *avahiHostname         // Hostname resolver
	units       map[string]*unit       // Discovered print/fax/scan units
}

// Delete deletes the service from avahiClient
func (service *avahiService) Delete() {
	clnt := service.clnt
	service.SetHostname(nil)

	service.svcResolver.Close()
	service.txtBrowser.Close()

	for name, un := range service.units {
		delete(service.units, name)
		un.Delete()
	}

	delete(clnt.services, service.key)
}

// SetPort sets service port
func (service *avahiService) SetPort(port uint16) {
	if service.port == port {
		return // Nothing changed
	}

	service.port = port
	for _, un := range service.units {
		un.SetPort(port)
	}
}

// SetHostname creates an association between service and hostname.
//
// If the new hostname is nil, service becomes without hostname
// association.
//
// Old association is cleaned up. If old hostname is not longer
// in use by some services, it will be deleted.
func (service *avahiService) SetHostname(hostname *avahiHostname) {
	if service.hostname == hostname {
		// Nothing changed
		return
	}

	if service.hostname != nil {
		service.hostname.services.Del(service)
		if service.hostname.services.Empty() {
			service.hostname.Delete()
		}
	}

	service.hostname = hostname
	if hostname != nil {
		hostname.services.Add(service)
	}
}

// AddUnit adds unit to the service
func (service *avahiService) AddUnit(name string, un *unit) {
	service.units[name] = un
	un.SetPort(service.port)

	if service.hostname != nil {
		service.hostname.addrs.ForEach(func(addr netip.Addr) {
			un.AddAddr(addr)
		})
	}
}

// DelUnit deletes unit from the service
func (service *avahiService) DelUnit(name string) {
	delete(service.units, name)
}

// GelUnit returns unit by the name (or nil if unit doesn't exist)
func (service *avahiService) GetUnit(name string) *unit {
	return service.units[name]
}

// addAddr called by avahiHostname.AddAddr to record newly
// resolved address.
func (service *avahiService) addAddr(addr netip.Addr) {
	for _, un := range service.units {
		un.AddAddr(addr)
	}
}

// delAddr called by avahiHostname.DelAddr to record newly
// resolved address.
func (service *avahiService) delAddr(addr netip.Addr) {
	for _, un := range service.units {
		un.DelAddr(addr)
	}
}

// avahiServiceKey identifies a particular instance of
// service of the particular type.
type avahiServiceKey struct {
	IfIdx        avahi.IfIndex  // Network interface index
	Proto        avahi.Protocol // Network protocol
	InstanceName string         // Service instance name
	SvcType      string         // Service type
	Domain       string         // Service domain
}

// FQDN returns the full-qualified domain name for the
// service instance.
func (key avahiServiceKey) FQDN() string {
	return avahi.DomainServiceNameJoin(key.InstanceName,
		key.SvcType, key.Domain)
}

// String returns string representation of the avahiServiceKey,
// for debugging.
func (key avahiServiceKey) String() string {
	return fmt.Sprintf("%q (%s%%%d)", key.FQDN(), key.Proto, key.IfIdx)
}

// HostnameKey makes avahiHostnameKey
func (key avahiServiceKey) HostnameKey(name string) avahiHostnameKey {
	return avahiHostnameKey{
		IfIdx:    key.IfIdx,
		Proto:    key.Proto,
		Hostname: name,
	}
}

// commonUnitID fills parts of discovery.UnitID, common for
// printers, scanners and faxout devices
func (key avahiServiceKey) commonUnitID(
	backend discovery.Backend) discovery.UnitID {

	var variant string

	switch key.Proto {
	case avahi.ProtocolIP4:
		variant = "ip4"
	case avahi.ProtocolIP6:
		variant = "ip6"
	}

	switch key.SvcType {
	case svcTypeIPP, svcTypeESCL:
		variant += "-http"
	case svcTypeIPPS, svcTypeESCLS:
		variant += "-https"
	}

	return discovery.UnitID{
		DNSSDName: key.InstanceName,
		Backend:   backend,
		Zone:      zone.Name(int(key.IfIdx)),
		Variant:   variant,
		SvcProto:  svcTypeToDiscoveryServiceProto(key.SvcType),
	}
}

// PrinterUnitID makes discovery.UnitID for printer
func (key avahiServiceKey) PrinterUnitID(backend discovery.Backend,
	txt txtPrinter) discovery.UnitID {

	id := key.commonUnitID(backend)

	id.UUID = txt.uuid
	id.Queue = txt.rp
	id.SvcType = discovery.ServicePrinter
	id.USBSerial = txt.usbSerial
	id.USBHWID = txt.usbHWID

	return id
}

// PrinterUnitID makes discovery.UnitID for fax
func (key avahiServiceKey) FaxoutUnitID(backend discovery.Backend,
	txt txtPrinter) discovery.UnitID {

	id := key.commonUnitID(backend)

	id.UUID = txt.uuid
	id.Queue = txt.rfo
	id.SvcType = discovery.ServiceFaxout
	id.USBSerial = txt.usbSerial
	id.USBHWID = txt.usbHWID

	return id
}

// ScannerUnitID makes discovery.UnitID for scanner
func (key avahiServiceKey) ScannerUnitID(backend discovery.Backend,
	txt txtScanner) discovery.UnitID {

	id := key.commonUnitID(backend)

	id.UUID = txt.uuid
	id.SvcType = discovery.ServiceScanner
	id.USBSerial = txt.usbSerial
	id.USBHWID = txt.usbHWID

	return id
}

// IsPrinter reports if service type is printer
func (key avahiServiceKey) IsPrinter() bool {
	return svcTypeIsPrinter(key.SvcType)
}

// avahiServiceKeyFromServiceBrowserEvent makes avahiServiceKey
// from the avahi.ServiceBrowserEvent.
func avahiServiceKeyFromServiceBrowserEvent(
	evnt *avahi.ServiceBrowserEvent) avahiServiceKey {

	return avahiServiceKey{
		IfIdx:        evnt.IfIdx,
		Proto:        evnt.Proto,
		InstanceName: evnt.InstanceName,
		SvcType:      evnt.SvcType,
		Domain:       evnt.Domain,
	}
}

// avahiServiceKeyFromResolverEvent makes avahiServiceKey
// from the avahi.ServiceResolverEvent.
func avahiServiceKeyFromResolverEvent(
	evnt *avahi.ServiceResolverEvent) avahiServiceKey {

	return avahiServiceKey{
		IfIdx:        evnt.IfIdx,
		Proto:        evnt.Proto,
		InstanceName: evnt.InstanceName,
		SvcType:      evnt.SvcType,
		Domain:       evnt.Domain,
	}
}

// avahiServiceKeyFromRecordBrowserEvent makes avahiServiceKey
// from the avahi.RecordBrowserEvent
func avahiServiceKeyFromRecordBrowserEvent(
	evnt *avahi.RecordBrowserEvent) avahiServiceKey {

	instance, svctype, domain := avahi.DomainServiceNameSplit(evnt.Name)

	return avahiServiceKey{
		IfIdx:        evnt.IfIdx,
		Proto:        evnt.Proto,
		InstanceName: instance,
		SvcType:      svctype,
		Domain:       domain,
	}
}

// avahiHostname is the per-hostname structure.
//
// It is created on demand, when some of discovered service
// instances contain reference to that hostname and manages
// resolving of IP addresses, associated with that hostname.
type avahiHostname struct {
	clnt        *avahiClient               // The owner
	key         avahiHostnameKey           // Identity
	addrBrowser *avahi.RecordBrowser       // A/AAAA record resolver
	addrs       generic.Set[netip.Addr]    // Resolved addresses
	services    generic.Set[*avahiService] // Dependent services
}

// Delete deletes the avahiHostname from avahiClient
func (hostname *avahiHostname) Delete() {
	clnt := hostname.clnt
	hostname.addrBrowser.Close()
	delete(clnt.hostnames, hostname.key)
}

// HasAddr reports if address already known
func (hostname *avahiHostname) HasAddr(addr netip.Addr) bool {
	return hostname.addrs.Contains(addr)
}

// addAddr adds the IP address, associated with the hostname.
func (hostname *avahiHostname) AddAddr(addr netip.Addr, ifidx avahi.IfIndex) {
	// Filter address according to the following rules
	//   - if service belongs to the loopback interface, only loopback
	//     addresses are allowed
	//   - link-local addresses must belong to the same interface as
	//     the service belongs
	//   - all other cases are allowed
	switch {
	case addr.IsLoopback():
		if hostname.key.IfIdx != loopback {
			return
		}

	case addr.IsLinkLocalUnicast():
		// In case Avahi will return all A/AAAA records for the
		// hostname, filter out link-local addresses that doesn't
		// belong to the querying interface.
		if ifidx != hostname.key.IfIdx {
			return
		}
	}

	// Add the address
	if !hostname.addrs.Contains(addr) {
		hostname.addrs.Add(addr)
		hostname.services.ForEach(func(service *avahiService) {
			service.addAddr(addr)
		})
	}
}

// delAddr deletes the IP address, that previously was associated
// with the hostname.
func (hostname *avahiHostname) DelAddr(addr netip.Addr) {
	if hostname.addrs.Contains(addr) {
		hostname.addrs.Del(addr)
		hostname.services.ForEach(func(service *avahiService) {
			service.delAddr(addr)
		})
	}
}

// avahiHostnameKey identifies a particular instance of
// the avahiHostname
type avahiHostnameKey struct {
	IfIdx    avahi.IfIndex  // Network interface index
	Proto    avahi.Protocol // Network protocol
	Hostname string         // Hostname string (FQDN)
}

// String returns string representation of the avahiHostnameKey,
// for debugging.
func (key avahiHostnameKey) String() string {
	return fmt.Sprintf("%q (%s%%%d)", key.Hostname, key.Proto, key.IfIdx)
}

// avahiHostnameKeyFromRecordBrowserEvent makes avahiHostnameKey
// from the avahi.RecordBrowserEvent
func avahiHostnameKeyFromRecordBrowserEvent(
	evnt *avahi.RecordBrowserEvent) avahiHostnameKey {

	return avahiHostnameKey{
		IfIdx:    evnt.IfIdx,
		Proto:    evnt.Proto,
		Hostname: evnt.Name,
	}
}

// avahiLookupFlagsTable contains mapping of LookupFlags to avahi.LookupFlags
var avahiLookupFlagsTable = [...]avahi.LookupFlags{
	0:               avahi.LookupUseMulticast,
	LookupClassical: avahi.LookupUseWideArea,
	LookupMulticast: avahi.LookupUseMulticast,
	LookupBoths:     avahi.LookupUseMulticast,
}

// avahiLookupFlags maps LookupFlags to avahi.LookupFlags
func avahiLookupFlags(flags LookupFlags) avahi.LookupFlags {
	return avahiLookupFlagsTable[flags&LookupBoths]
}
