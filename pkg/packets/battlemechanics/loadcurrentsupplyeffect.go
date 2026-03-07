package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.LoadSupplyEffectID

	return &LoadSupplyEffectPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadSupplyEffectID, "LoadSupplyEffect", func() packets.Packet {
		return NewLoadSupplyEffectPacket()
	})
}
