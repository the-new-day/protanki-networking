package garage

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.BuyKitID

	return &BuyKitPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BuyKitID, "BuyKit", func() packets.Packet {
		return NewBuyKitPacket()
	})
}
