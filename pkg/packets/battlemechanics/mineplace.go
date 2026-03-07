package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sent when a mine is placed.
type MinePlacePacket struct {
	packets.BasePacket
}

func NewMinePlacePacket() *MinePlacePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"mineId",
	}

	var id int32 = packets.MinePlaceID

	return &MinePlacePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MinePlaceID, "MinePlace", func() packets.Packet {
		return NewMinePlacePacket()
	})
}
