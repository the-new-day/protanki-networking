package garage

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Buy multiple items from garage, like supplies, xp boosts
type BuyMultipleItemsPacket struct {
	packets.BasePacket
}

func NewBuyMultipleItemsPacket() *BuyMultipleItemsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"item_id",
		"count",
		"base_cost",
	}

	var id int32 = packets.BuyMultipleItemsID

	return &BuyMultipleItemsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BuyMultipleItemsID, "BuyMultipleItems", func() packets.Packet {
		return NewBuyMultipleItemsPacket()
	})
}
