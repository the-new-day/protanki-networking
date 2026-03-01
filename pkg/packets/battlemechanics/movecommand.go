package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Receives movement data of a player from the server.
type MoveCommandPacket struct {
	packets.BasePacket
}

func NewMoveCommandPacket() *MoveCommandPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewMoveCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
	}

	attributes := []string{
		"username",
		"movement",
		"turretDirection",
	}

	var id int32 = 1516578027

	return &MoveCommandPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1516578027, "MoveCommand", func() packets.Packet {
		return NewMoveCommandPacket()
	})
}
