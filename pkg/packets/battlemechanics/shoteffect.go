package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Shot effect packet
type ShotEffectPacket struct {
	packets.BasePacket
}

func NewShotEffectPacket() *ShotEffectPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"effectId",
	}

	var id int32 = packets.ShotEffectID

	return &ShotEffectPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShotEffectID, "ShotEffect", func() packets.Packet {
		return NewShotEffectPacket()
	})
}
