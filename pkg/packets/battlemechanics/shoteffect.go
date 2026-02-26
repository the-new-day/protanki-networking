package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -1994318624

	return &ShotEffectPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
