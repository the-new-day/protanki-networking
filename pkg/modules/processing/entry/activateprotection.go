package entry

import (
	"github.com/the-new-day/probogo/pkg/modules/processing"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

type ActivateProtectionHandler struct {
	*processing.BasePacketHandler
	protection protection.Protection
}

func NewActivateProtectionHandler(base *processing.BasePacketHandler, protection protection.Protection) *ActivateProtectionHandler {
	return &ActivateProtectionHandler{
		BasePacketHandler: base,
		protection:        protection,
	}
}

func (h *ActivateProtectionHandler) Handle(packet packets.Packet) (handled packets.Packet, err error) {
	if packet.ID() != packets.ActivateProtectionID {
		return h.CallNext(packet)
	}

	keys := packet.Get("keys").([]byte)
	h.protection.Activate(keys)
	return packet, nil
}
