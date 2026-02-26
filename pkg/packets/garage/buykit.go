package garage

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Buy a kit
type BuyKitPacket struct {
	packets.BasePacket
}

func NewBuyKitPacket() *BuyKitPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"item_id",
		"base_cost",
	}

	var id int32 = -523392052

	return &BuyKitPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
