package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Purchasable Items
type LoadPurchasableItemsPacket struct {
	packets.BasePacket
}

func NewLoadPurchasableItemsPacket() *LoadPurchasableItemsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -300370823

	return &LoadPurchasableItemsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-300370823, "LoadPurchasableItemsPacket", func() packets.Packet {
		return NewLoadPurchasableItemsPacket()
	})
}
