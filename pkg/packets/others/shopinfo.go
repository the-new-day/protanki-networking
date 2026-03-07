package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Get shop info
type ShopInfoPacket struct {
	packets.BasePacket
}

func NewShopInfoPacket() *ShopInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"data",
	}

	var id int32 = packets.ShopInfoID

	return &ShopInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShopInfoID, "ShopInfo", func() packets.Packet {
		return NewShopInfoPacket()
	})
}
