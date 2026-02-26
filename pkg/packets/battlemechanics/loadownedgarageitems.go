package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -255516505

	return &LoadOwnedGarageItemsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
