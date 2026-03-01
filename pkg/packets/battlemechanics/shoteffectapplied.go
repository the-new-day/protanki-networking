package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Shot Effect Applied
type ShotEffectAppliedPacket struct {
	packets.BasePacket
}

func NewShotEffectAppliedPacket() *ShotEffectAppliedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
	}

	attributes := []string{
		"shooter",
		"hitPoint",
	}

	var id int32 = 546849203

	return &ShotEffectAppliedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(546849203, "ShotEffectApplied", func() packets.Packet {
		return NewShotEffectAppliedPacket()
	})
}
