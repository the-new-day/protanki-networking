package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player Start Position
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

	var id int32 = packets.PlayerStartPositionID

	return &PlayerStartPositionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PlayerStartPositionID, "PlayerStartPosition", func() packets.Packet {
		return NewPlayerStartPositionPacket()
	})
}
