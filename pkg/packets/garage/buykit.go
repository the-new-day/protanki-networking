package garage

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
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

func init() {
	packets.Register(-523392052, "BuyKitPacket", func() packets.Packet {
		return NewBuyKitPacket()
	})
}
