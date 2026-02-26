package shop

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Promocode applied successfully
type PromocodeSuccessPacket struct {
	packets.BasePacket
}

func NewPromocodeSuccessPacket() *PromocodeSuccessPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -1859441081

	return &PromocodeSuccessPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
