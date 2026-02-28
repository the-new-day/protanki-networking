package entry

import (
	"github.com/the-new-day/probogo/pkg/modules/processing"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

type SetCaptchaKeysHandler struct {
	*processing.BasePacketHandler
	protection protection.Protection
}

func NewSetCaptchaKeysHandler(base *processing.BasePacketHandler, protection protection.Protection) *SetCaptchaKeysHandler {
	return &SetCaptchaKeysHandler{
		BasePacketHandler: base,
		protection:        protection,
	}
}

func (h *SetCaptchaKeysHandler) Handle(packet packets.Packet) (handled packets.Packet, err error) {
	if packet.ID() != packets.SetCaptchaKeysID {
		return h.CallNext(packet)
	}

	// TODO: implement
	return packet, nil
}
