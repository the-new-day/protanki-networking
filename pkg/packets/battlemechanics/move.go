package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends your movement data to the server
type MovePacket struct {
	packets.BasePacket
}

func NewMovePacket() *MovePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(custom.NewMoveCodec()),
	}

	attributes := []string{
		"clientTime",
		"specificationID",
		"movement",
	}

	var id int32 = packets.MoveID

	return &MovePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MoveID, "Move", func() packets.Packet {
		return NewMovePacket()
	})
}
