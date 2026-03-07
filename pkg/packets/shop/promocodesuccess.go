package shop

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Promocode applied successfully
type PromocodeSuccessPacket struct {
	packets.BasePacket
}

func NewPromocodeSuccessPacket() *PromocodeSuccessPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.PromocodeSuccessID

	return &PromocodeSuccessPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PromocodeSuccessID, "PromocodeSuccess", func() packets.Packet {
		return NewPromocodeSuccessPacket()
	})
}
