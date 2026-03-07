package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Load Owned Garage Items
type LoadOwnedGarageItemsPacket struct {
	packets.BasePacket
}

func NewLoadOwnedGarageItemsPacket() *LoadOwnedGarageItemsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.LoadOwnedGarageItemsID

	return &LoadOwnedGarageItemsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadOwnedGarageItemsID, "LoadOwnedGarageItems", func() packets.Packet {
		return NewLoadOwnedGarageItemsPacket()
	})
}
