package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Use a supply in the battle
type UseSupplyPacket struct {
	packets.BasePacket
}

func NewUseSupplyPacket() *UseSupplyPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"supply_id",
	}

	var id int32 = -2102525054

	return &UseSupplyPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
