package entry

import (
	"github.com/the-new-day/probogo/pkg/modules/networking"
	"github.com/the-new-day/probogo/pkg/modules/processing"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

type EntryHandler struct {
	processing.BasePacketHandler
	protection protection.Protection
}

func NewEntryHandler(socket *networking.GameSocket, protection protection.Protection) processing.PacketHandler {
	handler := &EntryHandler{
		BasePacketHandler: *processing.NewBasePacketHandler(socket, protection),
		protection:        protection,
	}
	// TODO: add children
	return handler
}

func (h *EntryHandler) Handle(packet packets.Packet) (handled packets.Packet, err error) {
	return h.CallChild(packet)
}
