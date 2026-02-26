package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Flag has been taken
type FlagTakenPacket struct {
	packets.BasePacket
}

func NewFlagTakenPacket() *FlagTakenPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"flagteam",
	}

	var id int32 = -1282406496

	return &FlagTakenPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
