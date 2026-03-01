package shop

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Incorrect or expired promocode
type PromocodeFailedPacket struct {
	packets.BasePacket
}

func NewPromocodeFailedPacket() *PromocodeFailedPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -1850050333

	return &PromocodeFailedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1850050333, "PromocodeFailed", func() packets.Packet {
		return NewPromocodeFailedPacket()
	})
}
