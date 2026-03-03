package networking

import (
	"context"

	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

type PacketHandler struct {
	stream *PacketStream

	ctx    context.Context
	cancel context.CancelFunc

	inboundHandlers  []func(packets.Packet) packets.Packet
	outboundHandlers []func(packets.Packet) packets.Packet

	receiveErrorHandlers []func(PacketResult)
}

func NewPacketHandler(
	conn connection.Connection,
	protection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *PacketHandler {
	ctx, cancel := context.WithCancel(context.Background())

	return &PacketHandler{
		stream:               NewPacketStream(conn, protection, packetRegistry),
		ctx:                  ctx,
		cancel:               cancel,
		inboundHandlers:      []func(packets.Packet) packets.Packet{},
		outboundHandlers:     []func(packets.Packet) packets.Packet{},
		receiveErrorHandlers: []func(PacketResult){},
	}
}

func (ph *PacketHandler) OnInBound(handler func(packets.Packet) packets.Packet) {
	ph.inboundHandlers = append(ph.inboundHandlers, handler)
}

func (ph *PacketHandler) OnOutBound(handler func(packets.Packet) packets.Packet) {
	ph.outboundHandlers = append(ph.outboundHandlers, handler)
}

func (ph *PacketHandler) OnReceiveError(handler func(PacketResult)) {
	ph.receiveErrorHandlers = append(ph.receiveErrorHandlers, handler)
}

func (ph *PacketHandler) Send(packet packets.Packet) error {
	for _, handler := range ph.outboundHandlers {
		packet = handler(packet)
		if packet == nil {
			return nil // packet has been canceled
		}
	}
	return ph.stream.Send(packet)
}

func (ph *PacketHandler) SendWithoutNotify(packet packets.Packet) error {
	return ph.stream.Send(packet)
}

func (ph *PacketHandler) Run() {
	for res := range ph.stream.Packets(ph.ctx) {
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

func (ph *PacketHandler) Close() {
	ph.cancel()
}
