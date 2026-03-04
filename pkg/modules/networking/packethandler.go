package networking

import (
	"context"

	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

// PacketHandler is a wrapper for PacketStream that allows subscribing
// to inbound and outbound packets and sending packets.
// Subscribers can read and modify all packets before they are delivered.
type PacketHandler struct {
	stream *PacketStream

	inboundHandlers  []func(packets.Packet) packets.Packet
	outboundHandlers []func(packets.Packet) packets.Packet

	receiveErrorHandlers []func(PacketResult)
}

// NewPacketHandler creates PacketHandler instance.
// It creates an instance of the underlying PacketStream with provided connection, protection, and registry.
func NewPacketHandler(
	conn connection.Connection,
	protection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *PacketHandler {
	return &PacketHandler{
		stream:               NewPacketStream(conn, protection, packetRegistry),
		inboundHandlers:      []func(packets.Packet) packets.Packet{},
		outboundHandlers:     []func(packets.Packet) packets.Packet{},
		receiveErrorHandlers: []func(PacketResult){},
	}
}

// OnInBound adds a handler (subscriber) for inbound packets (coming from the connection).
// Handlers must return the packet, they can modify the packet by returning the modified version.
// If a nil is returned from the handler, the packet is considered canceled and is not delivered.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnInBound(handler func(packets.Packet) packets.Packet) {
	ph.inboundHandlers = append(ph.inboundHandlers, handler)
}

// OnOutBound adds a handler (subscriber) for outbound packets (coming to the connection).
// Handlers must return the packet, they can modify the packet by returning the modified version.
// If a nil is returned from the handler, the packet is considered canceled and is not delivered.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnOutBound(handler func(packets.Packet) packets.Packet) {
	ph.outboundHandlers = append(ph.outboundHandlers, handler)
}

// OnError adds a handler (subscriber) for error events.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnError(handler func(PacketResult)) {
	ph.receiveErrorHandlers = append(ph.receiveErrorHandlers, handler)
}

// Send sends the packet to the underlying PacketStream, notifying all outbound subscribers.
// If some subscriber returns nil, notifying stops and the packet is not deliverd (it's considered canceled).
func (ph *PacketHandler) Send(packet packets.Packet) error {
	for _, handler := range ph.outboundHandlers {
		packet = handler(packet)
		if packet == nil {
			return nil // packet has been canceled
		}
	}
	return ph.stream.Send(packet)
}

// Run starts execution. It reads inbound packets and notifies all inbound subscribers.
// If any error occurs, it notifies all error subscribers and continues to read.
// It is a blocking operation.
func (ph *PacketHandler) Run(ctx context.Context) {
	for res := range ph.stream.Inbound(ctx) {
		if res.Err != nil {
			for _, handler := range ph.receiveErrorHandlers {
				handler(res)
			}
			continue
		}

		for _, handler := range ph.inboundHandlers {
			res.Packet = handler(res.Packet)
			if res.Packet == nil {
				break
			}
		}
	}
}

// ActivateProtection calls Activate(keys) on the underlying Protection instance.
func (ph *PacketHandler) ActivateProtection(keys []byte) {
	ph.stream.ActivateProtection(keys)
}
