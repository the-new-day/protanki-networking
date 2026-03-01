package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Change location in shop
type BuyFromShopPacket struct {
	packets.BasePacket
}

func NewBuyFromShopPacket() *BuyFromShopPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"itemId",
		"itemType",
	}

	var id int32 = 880756819

	return &BuyFromShopPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(880756819, "BuyFromShop", func() packets.Packet {
		return NewBuyFromShopPacket()
	})
}
