package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -624217047

	return &MinePlacePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-624217047, "MinePlacePacket", func() packets.Packet {
		return NewMinePlacePacket()
	})
}
