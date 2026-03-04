package proxy

import (
	"context"
	"log"

	"github.com/the-new-day/probogo/pkg/modules/networking"
	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
	"github.com/the-new-day/probogo/pkg/packets/network"
)

// Proxy is a universal MITM proxy for game-client communications.
// It allows intercepting and modifying packets.
type Proxy struct {
	serverHandler *networking.PacketHandler
	clientHandler *networking.PacketHandler
}

func NewProxy(
	serverConn connection.Connection,
	serverProtection protection.Protection,
	clientConn connection.Connection,
	clientProtection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *Proxy {
	p := &Proxy{
		serverHandler: networking.NewPacketHandler(serverConn, serverProtection, packetRegistry),
		clientHandler: networking.NewPacketHandler(clientConn, clientProtection, packetRegistry),
	}

	p.serverHandler.OnInBound(p.handleActivateProtection) // always being called first

	// TODO: remove logs
	p.OnClientToServer(func(packet packets.Packet) packets.Packet {
		template := "[CLIENT]: [ID: %d | Name: %s | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), packets.GetName(packet.ID()), 8+len(packet.EncryptedRawData()), packet.EncryptedRawData())

		p.serverHandler.Send(packet)
		return packet
	})

	p.OnServerToClient(func(packet packets.Packet) packets.Packet {
		template := "[SERVER]: [ID: %d | Name: %s | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), packets.GetName(packet.ID()), 8+len(packet.EncryptedRawData()), packet.EncryptedRawData())

		p.clientHandler.Send(packet)
		return packet
	})
	return p
}

func (p *Proxy) OnServerToClient(handler func(packets.Packet) packets.Packet) {
	p.serverHandler.OnInBound(handler)
}

func (p *Proxy) OnClientToServer(handler func(packets.Packet) packets.Packet) {
	p.clientHandler.OnInBound(handler)
}

func (p *Proxy) OnError(handler func(networking.PacketResult)) {
	p.serverHandler.OnError(handler)
	p.clientHandler.OnError(handler)
}

func (p *Proxy) Run(ctx context.Context) {
	go p.serverHandler.Run(ctx)
	go p.clientHandler.Run(ctx)
	<-ctx.Done()
}

// handleActivateProtection is a handler for the ActivateProtection packet.
// It activates server protection with provided keys and sends protection keys to the client.
func (p *Proxy) handleActivateProtection(packet packets.Packet) packets.Packet {
	if activateProt, ok := packet.(*network.ActivateProtectionPacket); ok {
		template := "[SERVER]: [ID: %d | Name: %s | Len: %d | RawData: % x]"
		log.Printf(template, packet.ID(), packets.GetName(packet.ID()), 8+len(packet.EncryptedRawData()), packet.EncryptedRawData())

		keys := packets.Attr[[]byte]("keys", activateProt)
		p.serverHandler.ActivateProtection(keys)

		p.clientHandler.Send(activateProt)
		p.clientHandler.ActivateProtection(keys)

		return nil // TODO: add listeners, it can't be passed further (messes up the state of Protection)
	}

	return packet
}
