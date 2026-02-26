package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Receives movement data of a player from the server.
type MovedPacket struct {
	packets.BasePacket
}

func NewMovedPacket() *MovedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewMoveCodec()),
	}

	attributes := []string{
		"username",
		"movement",
	}

	var id int32 = -64696933

	return &MovedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
