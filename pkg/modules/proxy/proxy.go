package proxy

import (
	"context"

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

	clientProtectionKeys []byte
}

func NewProxy(
	serverConn connection.Connection,
	serverProtection protection.Protection,
	clientConn connection.Connection,
	clientProtection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *Proxy {
	p := &Proxy{
		serverHandler:        networking.NewPacketHandler(serverConn, serverProtection, packetRegistry),
		clientHandler:        networking.NewPacketHandler(clientConn, clientProtection, packetRegistry),
		clientProtectionKeys: []byte{},
	}

	p.serverHandler.OnInBound(p.handleActivateProtection) // always being called first
	return p
}

func (p *Proxy) OnServerToClient(handler func(packets.Packet) packets.Packet) {
	p.serverHandler.OnInBound(handler)
}

func (p *Proxy) OnClientToServer(handler func(packets.Packet) packets.Packet) {
	p.clientHandler.OnInBound(handler)
}

func (p *Proxy) Run(ctx context.Context) {
	p.clientHandler.ActivateProtection(p.clientProtectionKeys)

	go p.clientHandler.Run(ctx)
	go p.serverHandler.Run(ctx)
}

// SetClientProtectionKeys sets protection keys used to activate protection during Run().
// Default: empty set of keys.
// It should be called before Run().
func (p *Proxy) SetClientProtectionKeys(keys []byte) {
	buf := make([]byte, len(keys))
	copy(buf, keys)
	p.clientProtectionKeys = buf
}

// handleActivateProtection is a handler for the ActivateProtection packet.
// It activates server protection with provided keys and sets "fake" client protection keys.
func (p *Proxy) handleActivateProtection(packet packets.Packet) packets.Packet {
	if activateProt, ok := packet.(*network.ActivateProtectionPacket); ok {
		keys := packets.Attr[[]byte]("keys", activateProt)
		p.serverHandler.ActivateProtection(keys)
		activateProt.Set("keys", p.clientProtectionKeys) // sending "fake" keys to the client
	}

	return packet
}
