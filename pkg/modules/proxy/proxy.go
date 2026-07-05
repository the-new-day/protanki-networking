package proxy

import (
	"bytes"
	"context"
	"sync"

	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/modules/networking"
	"github.com/the-new-day/protanki-networking/pkg/modules/networking/connection"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
	"github.com/the-new-day/protanki-networking/pkg/packets"
	"github.com/the-new-day/protanki-networking/pkg/packets/network"
)

var emptyOnActivateProtection = func(keys []byte) {}

// Proxy is a universal MITM proxy for game-client communications.
// It allows intercepting and modifying packets,
// and transfers (possibly modified) packets between the server and the client.
//
// It accepts two instances of Protection: one with flipDirection = true, another with flipDirection = false.
// One Protection imitates Protection of the server, another imitates client Protection.
// Packets from the server are intercepted and decrypted using "client" Protection,
// then they are encrypted again with the "server" Protection.
//
// This allows reading, modifying, or imitating packets.
type Proxy struct {
	serverHandler *PacketHandler
	clientHandler *PacketHandler

	onActivateProtection []func(packet packets.Packet)
}

// NewProxy creates new Proxy instance.
// It also registers InBound subscriber for handling ActivateProtection packet,
// which sets up the protection keys for the server and client communication
// and prevents it from being passed to other packets (they are not allowed modify it).
func NewProxy(
	serverConn connection.Connection,
	serverProtection protection.Protection,
	clientConn connection.Connection,
	clientProtection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *Proxy {
	p := &Proxy{}

	p.serverHandler = NewPacketHandler(serverConn, serverProtection, packetRegistry, p.activateProtection)
	p.clientHandler = NewPacketHandler(clientConn, clientProtection, packetRegistry, emptyOnActivateProtection)

	p.serverHandler.OnInBound(p.notifyActivateProtection) // always being called first, cancels the packet

	p.clientHandler.OnInBoundFinal(func(packet packets.Packet) {
		p.SendToServer(packet)
	})

	p.serverHandler.OnInBoundFinal(func(packet packets.Packet) {
		p.SendToClient(packet)
	})

	p.clientHandler.OnError(func(pr networking.PacketResult) {
		handleError(p.serverHandler, pr)
	})

	p.serverHandler.OnError(func(pr networking.PacketResult) {
		handleError(p.clientHandler, pr)
	})
	return p
}

// OnServerToClient adds a subscriber for packets coming from the server.
// Subscribers can modify the packet by returning the edited version,
// and can cancel if from delivering by returning nil
// (notifying will also stop, but all handlers that already received the packet won't be notified that the packet is canceled).
//
// ActivateProtection packet won't be passed to any handler registered here.
// Use OnActivateProtection for listening (modifying is not allowed).
// All subscribers are notified in the order of adding.
func (p *Proxy) OnServerToClient(handler func(packets.Packet) packets.Packet) {
	p.serverHandler.OnInBound(handler)
}

// OnClientToServer adds a subscriber for packets coming from the client.
// Subscribers can modify the packet by returning the edited version,
// and can cancel if from delivering by returning nil
// (notifying will also stop, but all handlers that already received the packet won't be notified that the packet is canceled).
// All subscribers are notified in the order of adding.
func (p *Proxy) OnClientToServer(handler func(packets.Packet) packets.Packet) {
	p.clientHandler.OnInBound(handler)
}

// OnClientError adds a subscriber for errors happened during communication with the client.
// All subscribers are notified in the order of adding.
func (p *Proxy) OnClientError(handler func(networking.PacketResult)) {
	p.clientHandler.OnError(handler)
}

// OnServerError adds a subscriber for errors happened during communication with the server.
// All subscribers are notified in the order of adding.
func (p *Proxy) OnServerError(handler func(networking.PacketResult)) {
	p.serverHandler.OnError(handler)
}

// OnActivateProtection adds a subscriber for ActivateProtection packet coming from the server.
// Subscribers can't modify the packet (they receive only a copy of the packet).
// All subscribers are notified in the order of adding.
func (p *Proxy) OnActivateProtection(listener func(packets.Packet)) {
	p.onActivateProtection = append(p.onActivateProtection, listener)
}

// Run starts the proxy. It runs both client and server communications.
// It is a blocking operation, it can be stopped using the ctx.
func (p *Proxy) Run(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		p.serverHandler.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		p.clientHandler.Run(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
}

// SendToClient sends packet to the client.
// Subscribers are not notified.
func (p *Proxy) SendToClient(packet packets.Packet) error {
	return p.clientHandler.Send(packet)
}

// SendToServer sends packet to the server.
// Subscribers are not notified.
func (p *Proxy) SendToServer(packet packets.Packet) error {
	return p.serverHandler.Send(packet)
}

// notifyActivateProtection reads ActivateProtection packet,
// notifies all onActivateProtection subscribers and returns nil (cancels the packet).
// The protection gets activated by activateProtection() callback for PacketStream.
func (p *Proxy) notifyActivateProtection(packet packets.Packet) packets.Packet {
	if activateProt, ok := packet.(*network.ActivateProtectionPacket); ok {
		for _, listener := range p.onActivateProtection {
			listener(packets.Clone(&activateProt.BasePacket))
		}
		return nil
	}

	return packet
}

// activateProtection activates client and server handlers' protection.
// It also sends ActivateProtection directly to the real client (using clientHandler).
func (p *Proxy) activateProtection(keys []byte) {
	packet := network.NewActivateProtectionPacketWithKeys(keys)
	data, _ := packet.Wrap(&protection.EmptyProtection{})

	p.clientHandler.SendRaw(data.Bytes()) // writing directly to avoid changing the state of protection

	p.clientHandler.ActivateProtection(keys)
	p.serverHandler.ActivateProtection(keys)
}

// handleError sends packet to the handler.
// If the packet was compressed, sets the compression bit.
func handleError(packetHandler *PacketHandler, pr networking.PacketResult) {
	buf := &bytes.Buffer{}
	intCodec := &primitive.IntCodec{}

	if pr.WasCompressed {
		pr.Length |= 0x40000000
	}

	intCodec.Encode(pr.Length, buf)
	intCodec.Encode(pr.ID, buf)
	packetHandler.SendRaw(buf.Bytes())

	packetHandler.SendRawEncrypted(pr.Data)
}
