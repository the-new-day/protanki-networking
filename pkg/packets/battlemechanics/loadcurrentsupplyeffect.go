package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Current Supply Effect
type LoadSupplyEffectPacket struct {
	packets.BasePacket
}

func NewLoadSupplyEffectPacket() *LoadSupplyEffectPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = 417965410

	return &LoadSupplyEffectPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(417965410, "LoadSupplyEffectPacket", func() packets.Packet {
		return NewLoadSupplyEffectPacket()
	})
}
