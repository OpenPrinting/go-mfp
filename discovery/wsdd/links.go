// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// UDP data links

package wsdd

import (
	"context"
	"net/netip"
	"sync"

	"github.com/alexpevzner/mfp/discovery/netstate"
	"github.com/alexpevzner/mfp/internal/generic"
	"github.com/alexpevzner/mfp/uuid"
	"github.com/alexpevzner/mfp/wsd"
)

// links dynamically manages per-local-address UDP links.
type links struct {
	back   *backend                           // Parent backend
	netmon *netstate.Notifier                 // Network state monitor
	mconn4 *mconn                             // For recv of IP4 multicasts
	mconn6 *mconn                             // For recv of IP6 multicasts
	table  map[netip.Addr]*link               // Per-local address links
	lock   sync.Mutex                         // links.table lock
	ports  *generic.LockedSet[netip.AddrPort] // Set of Local ports

	// querier.procNetmon closing synchronization
	ctxNetmon    context.Context    // Cancelable context for procNetmon
	cancelNetmon context.CancelFunc // Its cancellation function
	doneNetmon   sync.WaitGroup     // Wait for procNetmon termination

	// querier.procMconn closing synchronization
	doneMconn sync.WaitGroup // Wait for procMconn termination
}

// newLinks creates a new links structure
func newLinks(back *backend) (*links, error) {
	// Create multicast sockets
	mconn4, err := newMconn(wsddMulticastIP4)
	if err != nil {
		return nil, err
	}

	mconn6, err := newMconn(wsddMulticastIP6)
	if err != nil {
		mconn4.Close()
		return nil, err
	}

	// Create links structure
	lt := &links{
		back:   back,
		netmon: netstate.NewNotifier(),
		mconn4: mconn4,
		mconn6: mconn6,
		table:  make(map[netip.Addr]*link),
		ports:  generic.NewLockedSet[netip.AddrPort](),
	}

	return lt, nil
}

// Start starts links operations.
func (lt *links) Start() {
	// Start links.procNetmon
	lt.ctxNetmon, lt.cancelNetmon = context.WithCancel(lt.back.ctx)
	lt.doneNetmon.Add(1)
	go lt.procNetmon()

	// Start links.procMconn, one per connection
	lt.doneMconn.Add(2)
	go lt.procMconn(lt.mconn4)
	go lt.procMconn(lt.mconn6)
}

// Close closes links table and all links it owns.
func (lt *links) Close() {
	// Stop procNetmon
	lt.cancelNetmon()
	lt.doneNetmon.Wait()

	// Stop multicasts reception
	lt.mconn4.Close()
	lt.mconn6.Close()
	lt.doneMconn.Wait()

	// Close each individual link
	lt.lock.Lock()

	lt.ports.Clear()
	for addr, l := range lt.table {
		l.Close()
		delete(lt.table, addr)
	}

	lt.lock.Unlock()
}

// Add adds a local address and corresponding link
func (lt *links) add(addr netstate.Addr) {
	// Ignore non-multicast links
	flags := addr.Interface().Flags()
	if !flags.All(netstate.NetIfMulticast) {
		return
	}

	// Add link
	l := newLink(lt, addr)

	lt.lock.Lock()
	lt.table[addr.Addr()] = l
	lt.lock.Unlock()
}

// Del deletes local address
func (lt *links) del(addr netstate.Addr) {
	// Ignore non-multicast links
	flags := addr.Interface().Flags()
	if !flags.All(netstate.NetIfMulticast) {
		return
	}

	// Del link
	lt.lock.Lock()
	l := lt.table[addr.Addr()]
	delete(lt.table, addr.Addr())
	lt.lock.Unlock()

	l.Close()
}

// IsLocalPort reports if given port belongs to our local ports
func (lt *links) IsLocalPort(addr netip.AddrPort) bool {
	return lt.ports.Contains(addr)
}

// netmonproc processes netstate.Notifier events.
// It runs on its own goroutine.
func (lt *links) procNetmon() {
	defer lt.doneNetmon.Done()

	for {
		evnt, err := lt.netmon.Get(lt.ctxNetmon)
		if err != nil {
			return
		}

		lt.back.debug("%s", evnt)

		switch evnt := evnt.(type) {
		case netstate.EventAddPrimaryAddress:
			lt.add(evnt.Addr)
		case netstate.EventDelPrimaryAddress:
			lt.del(evnt.Addr)
		}
	}
}

// procMconn receives UDP multicast messages from the multicast conection.
func (lt *links) procMconn(mc *mconn) {
	defer lt.doneMconn.Done()

	for {
		var buf [65536]byte
		n, from, cmsg, err := mc.RecvFrom(buf[:])

		if mc.IsClosed() {
			return
		}

		if err != nil {
			lt.back.error("UDP recv: %s", err)
			continue
		}

		lt.back.input(buf[:n], from, mc.LocalAddrPort(), cmsg.IfIndex)
	}
}

// link is a per-local address link. It implements sending
// of the UDP multicast packets from that address and reception
// of responses.
//
// It also automatically sends Probe requests.
type link struct {
	parent     *links         // Parent links table
	addr       netstate.Addr  // Local address
	probeMsg   []byte         // Probe message
	probeSched *sched         // Probe scheduler
	dest       netip.AddrPort // Destination (multicast) address
	conn       *uconn         // Connection for sending UDP multicasts
	doneProber sync.WaitGroup // Wait for procProber termination
	doneReader sync.WaitGroup // Wait for procReader termination
}

// newLink creates a new link
func newLink(lt *links, addr netstate.Addr) *link {
	l := &link{
		addr:       addr,
		parent:     lt,
		probeSched: newSched(false),
	}

	if addr.Is4() {
		l.dest = wsddMulticastIP4
	} else {
		l.dest = wsddMulticastIP6
	}

	l.doneProber.Add(1)
	go l.procProber()

	return l
}

// close closes the link
func (l *link) Close() {
	// Kill link.procProber
	l.probeSched.Close()
	l.doneProber.Wait()

	// Kill l.procReader, if it is started
	if l.conn != nil {
		l.parent.ports.Del(l.conn.LocalAddrPort())
		l.conn.Close()
		l.doneReader.Wait()
	}
}

// procProber creates UDP connection on demand and sends probes.
// It runs on its own goroutine and paced by the sa.probeSched scheduler.
func (l *link) procProber() {
	defer l.doneProber.Done()

	back := l.parent.back

	for {
		evnt := <-l.probeSched.Chan()
		switch evnt {
		case schedClosed:
			return

		case schedNewMessage:
			// Open connection on demand
			if l.conn == nil {
				var err error
				l.conn, err = newUconn(l.addr, 0)
				if err != nil {
					back.debug("%s", err)
				}

				if l.conn != nil {
					l.parent.ports.Add(
						l.conn.LocalAddrPort())
					l.doneReader.Add(1)
					go l.procReader()
				}
			}

			// Update l.probeMsg
			l.updateProbeMsg()

		case schedSend:
			if l.conn != nil {
				l.conn.WriteToUDPAddrPort(l.probeMsg, l.dest)
				back.debug("%s message sent to %s%%%s",
					wsd.ActProbe, l.dest,
					l.addr.Interface().Name())
			}
		}
	}
}

// procReader runs on its own goroutine and receives messages from the l.conn.
func (l *link) procReader() {
	defer l.doneReader.Done()

	ifidx := l.addr.Interface().Index()
	to := l.conn.LocalAddrPort()
	back := l.parent.back

	for {
		// Receive next packet
		var buf [65536]byte
		n, from, err := l.conn.RecvFrom(buf[:])

		if l.conn.IsClosed() {
			return
		}

		if err != nil {
			l.parent.back.error("UDP recv: %s", err)
			continue
		}

		// Dispatch the packet
		back.input(buf[:n], from, to, ifidx)
	}
}

// updateProbeMsg updates l.probeMsg
func (l *link) updateProbeMsg() {
	msgid := wsd.AnyURI(uuid.Must(uuid.Random()).URN())
	msg := wsd.Msg{
		Header: wsd.Header{
			Action:    wsd.ActProbe,
			MessageID: msgid,
			To:        wsd.ToDiscovery,
		},
		Body: wsd.Probe{
			Types: wsd.TypeDevice,
		},
	}
	l.probeMsg = msg.Encode()
}
