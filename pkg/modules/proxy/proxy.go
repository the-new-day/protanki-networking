package proxy

import (
	"bytes"
	"context"
	"log"

	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/modules/networking"
	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
	"github.com/the-new-day/probogo/pkg/packets/network"
)

// Proxy is a universal MITM proxy for game-client communications.
// It allows intercepting and modifying packets,
// and transfers (possibly modified) packets between the server and the client.
//
// It accepts two instances of Protection: one with flipDirection = true, another with flipDirection = false.
// One Protection imitates Protection of the server, another imitates client Protection.
// Packets from the server are intercepted and decrypted using "client" Protection,
// then they are encrypted again with the "server" Protection.
type Proxy struct {
	serverHandler *networking.PacketHandler
	clientHandler *networking.PacketHandler

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

	p.serverHandler = networking.NewPacketHandler(serverConn, serverProtection, packetRegistry, p.activateProtection)
	p.clientHandler = networking.NewPacketHandler(clientConn, clientProtection, packetRegistry, func([]byte) {})

	p.serverHandler.OnInBound(p.handleActivateProtection) // always being called first

	// TODO: remove logs
	// TODO: make those execute the last

	p.OnClientToServer(func(packet packets.Packet) packets.Packet {
		template := "[CLIENT]: [ID: %d | Name: %s | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), packets.GetName(packet.ID()), 8+len(packet.RawData()), packet.RawData())

		p.serverHandler.Send(packet)
		return packet
	})

	p.OnServerToClient(func(packet packets.Packet) packets.Packet {
		template := "[SERVER]: [ID: %d | Name: %s | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), packets.GetName(packet.ID()), 8+len(packet.RawData()), packet.RawData())

		p.clientHandler.Send(packet)
		return packet
	})

	p.OnClientError(func(pr networking.PacketResult) {
		log.Printf("[ERROR | CLIENT]: ID: %d | Len: %d | Packet: % x | %v", pr.ID, pr.Length, pr.RawHex, pr.Err)
		buf := &bytes.Buffer{}
		intCodec := &primitive.IntCodec{}

		intCodec.Encode(pr.Length, buf)
		p.serverHandler.SendRaw(buf.Bytes())
		buf.Reset()

		intCodec.Encode(pr.ID, buf)
		p.serverHandler.SendRaw(buf.Bytes())

		p.serverHandler.SendRawEncrypted(pr.RawHex)
	})

	p.OnServerError(func(pr networking.PacketResult) {
		log.Printf("[ERROR | SERVER]: ID: %d | Len: %d | Packet: % x | %v", pr.ID, pr.Length, pr.RawHex, pr.Err)
		buf := &bytes.Buffer{}
		intCodec := &primitive.IntCodec{}

		intCodec.Encode(pr.Length, buf)
		p.clientHandler.SendRaw(buf.Bytes())
		buf.Reset()

		intCodec.Encode(pr.ID, buf)
		p.clientHandler.SendRaw(buf.Bytes())

		p.clientHandler.SendRawEncrypted(pr.RawHex)
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

// OnActivateProtection adds a subscriber for ActivateProtection packet coming from the server.
// Subscribers can't modify the packet (they receive only a copy of the packet).
// All subscribers are notified in the order of adding.
func (p *Proxy) OnActivateProtection(listener func(packets.Packet)) {
	p.onActivateProtection = append(p.onActivateProtection, listener)
}

// OnClientError adds a subscriber for errors happened during communication with the client.
// All subscribers are notified in the order of adding.
func (p *Proxy) OnClientError(handler func(networking.PacketResult)) {
	p.clientHandler.OnError(handler)
}

// OnClientError adds a subscriber for errors happened during communication with the server.
// All subscribers are notified in the order of adding.
func (p *Proxy) OnServerError(handler func(networking.PacketResult)) {
	p.serverHandler.OnError(handler)
}

// Run starts the proxy. It runs both client and server communications.
// It is a blocking operation, it can be stopped using the ctx.
func (p *Proxy) Run(ctx context.Context) {
	go p.serverHandler.Run(ctx)
	go p.clientHandler.Run(ctx)

	<-ctx.Done()
}

func (p *Proxy) handleActivateProtection(packet packets.Packet) packets.Packet {
	if activateProt, ok := packet.(*network.ActivateProtectionPacket); ok {
		template := "[ActivateProtection]: [ID: %d | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), 8+len(packet.RawData()), packet.RawData())

		for _, listener := range p.onActivateProtection {
			listener(packets.Clone(&activateProt.BasePacket))
		}
		return nil
	}

	return packet
}

func (p *Proxy) activateProtection(keys []byte) {
	packet := network.NewActivateProtectionPacketWithKeys(keys)
	data, _ := packet.Wrap(&protection.EmptyProtection{})

	p.clientHandler.SendRaw(data.Bytes()) // writing directly to avoid changing the state of protection

	p.clientHandler.ActivateProtection(keys)
	p.serverHandler.ActivateProtection(keys)
}
