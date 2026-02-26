package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Effect Aftermath
type EffectAftermathPacket struct {
	packets.BasePacket
}

func NewEffectAftermathPacket() *EffectAftermathPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.ByteCodec{}),
	}

	attributes := []string{
		"username",
		"effectId",
		"duration",
		"activeAfterDeath",
		"effectLevel",
	}

	var id int32 = -1639713644

	return &EffectAftermathPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
