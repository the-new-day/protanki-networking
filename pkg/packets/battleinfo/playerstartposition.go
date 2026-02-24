package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player start position.
type PlayerStartPositionPacket struct {
	packets.BasePacket
}

func NewPlayerStartPositionPacket() *PlayerStartPositionPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
	}
	attributes := []string{
		"position",
		"orientation",
	}

	var id int32 = -157204477

	return &PlayerStartPositionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
