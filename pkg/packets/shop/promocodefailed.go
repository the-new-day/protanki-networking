package shop

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Incorrect or expired promocode
type PromocodeFailedPacket struct {
	packets.BasePacket
}

func NewPromocodeFailedPacket() *PromocodeFailedPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.PromocodeFailedID

	return &PromocodeFailedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PromocodeFailedID, "PromocodeFailed", func() packets.Packet {
		return NewPromocodeFailedPacket()
	})
}
