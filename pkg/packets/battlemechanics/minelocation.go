package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sent when a mine is placed or removed.
type MineLocationPacket struct {
	packets.BasePacket
}

func NewMineLocationPacket() *MineLocationPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"mineId",
		"x",
		"y",
		"z",
		"userId",
	}

	var id int32 = 272183855

	return &MineLocationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
