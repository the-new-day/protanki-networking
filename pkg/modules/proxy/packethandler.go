package proxy

import (
	"context"

	"github.com/the-new-day/protanki-networking/pkg/modules/networking"
	"github.com/the-new-day/protanki-networking/pkg/modules/networking/connection"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// PacketHandler is a wrapper for PacketStream that allows subscribing
// to inbound and outbound packets and sending packets.
// Subscribers can read and modify all packets before they are delivered.
type PacketHandler struct {
	stream *networking.PacketStream

	onActivateProtection func([]byte)

	inboundHandlers  []func(packets.Packet) packets.Packet
	outboundHandlers []func(packets.Packet) packets.Packet

	inboundFinalHandlers  []func(packets.Packet)
	outboundFinalHandlers []func(packets.Packet)

	receiveErrorHandlers []func(networking.PacketResult)
}

// NewPacketHandler creates PacketHandler instance.
// It creates an instance of the underlying PacketStream with provided connection, protection, and registry.
// onActivateProtection is called when the ActivateProtection packet is received from the server,
// and it should be used to activate the protection (it is not activated by the PacketHandler).
func NewPacketHandler(
	conn connection.Connection,
	protection protection.Protection,
	packetRegistry *packets.PacketRegistry,
	onActivateProtection func([]byte),
) *PacketHandler {
	return &PacketHandler{
		stream:               networking.NewPacketStream(conn, protection, packetRegistry, onActivateProtection),
		inboundHandlers:      []func(packets.Packet) packets.Packet{},
		outboundHandlers:     []func(packets.Packet) packets.Packet{},
		receiveErrorHandlers: []func(networking.PacketResult){},
	}
}

// OnInBound adds a handler (subscriber) for inbound packets (coming from the connection).
// Handlers must return the packet, they can modify the packet by returning the modified version.
// If a nil is returned from the handler, the notifying stops.
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

// OnOutBoundFinal adds a final handler (subscriber) for outbound packets (coming to the connection).
// Final handlers are called after regular outbound handlers, and they can't modify or cancel packets.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnOutBoundFinal(handler func(packets.Packet)) {
	ph.outboundFinalHandlers = append(ph.outboundFinalHandlers, handler)
}

// OnInBoundFinal adds a final handler (subscriber) for inbound packets (coming from the connection).
// Final handlers are called after regular inbound handlers, and they can't modify or cancel packets.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnInBoundFinal(handler func(packets.Packet)) {
	ph.inboundFinalHandlers = append(ph.inboundFinalHandlers, handler)
}

// OnError adds a handler (subscriber) for error events.
// All handlers are notified in the order of adding.
func (ph *PacketHandler) OnError(handler func(networking.PacketResult)) {
	ph.receiveErrorHandlers = append(ph.receiveErrorHandlers, handler)
}

// Send sends the packet to the underlying PacketStream, notifying all outbound subscribers.
// If some subscriber returns nil, notifying stops and the packet is not delivered (it's considered canceled).
func (ph *PacketHandler) Send(packet packets.Packet) error {
	packet = ph.notifySend(packet)
	if packet == nil {
		return nil
	}
	return ph.stream.Send(packet)
}

// SendRaw sends raw data without encryption.
// Subscribers are not notified.
func (ph *PacketHandler) SendRaw(data []byte) error {
	return ph.stream.SendRaw(data)
}

// SendRaw encrypts and sends raw data.
// Subscribers are not notified.
func (ph *PacketHandler) SendRawEncrypted(data []byte) error {
	return ph.stream.SendRawEncrypted(data)
}

// Run starts execution. It reads inbound packets and notifies all inbound subscribers.
// If any error occurs, it notifies all error subscribers and continues to read.
// It is a blocking operation.
func (ph *PacketHandler) Run(ctx context.Context) {
	for res := range ph.stream.Inbound(ctx) {
		if res.Err != nil {
			ph.notifyError(res)
			continue
		}

		ph.notifyInbound(res.Packet)
	}
}

func (ph *PacketHandler) notifySend(packet packets.Packet) packets.Packet {
	for _, handler := range ph.outboundHandlers {
		if packet == nil {
			return nil // packet has been canceled
		}
		packet = handler(packet)
	}
	if packet != nil {
		for _, handler := range ph.outboundFinalHandlers {
			handler(packet)
		}
	}
	return packet
}

func (ph *PacketHandler) notifyError(packetRes networking.PacketResult) {
	for _, handler := range ph.receiveErrorHandlers {
		handler(packetRes)
	}
}

func (ph *PacketHandler) notifyInbound(packet packets.Packet) {
	for _, handler := range ph.inboundHandlers {
		if packet == nil {
			break
		}
		packet = handler(packet)
	}
	if packet != nil {
		for _, handler := range ph.inboundFinalHandlers {
			handler(packet)
		}
	}
}

// Activate protection calls Activate(keys) on the underlying protection.
func (ph *PacketHandler) ActivateProtection(keys []byte) {
	ph.stream.ActivateProtection(keys)
}
