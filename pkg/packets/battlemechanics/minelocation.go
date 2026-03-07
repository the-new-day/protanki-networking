package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.MineLocationID

	return &MineLocationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MineLocationID, "MineLocation", func() packets.Packet {
		return NewMineLocationPacket()
	})
}
