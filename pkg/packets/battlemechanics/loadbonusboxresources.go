package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Bonus Box Resources
type LoadBonusBoxResourcesPacket struct {
	packets.BasePacket
}

func NewLoadBonusBoxResourcesPacket() *LoadBonusBoxResourcesPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = 228171466

	return &LoadBonusBoxResourcesPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
