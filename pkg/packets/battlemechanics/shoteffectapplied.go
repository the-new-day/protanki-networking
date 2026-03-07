package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ShotEffectAppliedID

	return &ShotEffectAppliedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShotEffectAppliedID, "ShotEffectApplied", func() packets.Packet {
		return NewShotEffectAppliedPacket()
	})
}
