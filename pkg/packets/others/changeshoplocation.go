package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.BuyFromShopID

	return &BuyFromShopPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BuyFromShopID, "BuyFromShop", func() packets.Packet {
		return NewBuyFromShopPacket()
	})
}
