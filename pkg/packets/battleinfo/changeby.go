package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Change By
type ChangeByPacket struct {
	packets.BasePacket
}

func NewChangeByPacket() *ChangeByPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"changeBy",
	}

	var id int32 = packets.ChangeByID

	return &ChangeByPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ChangeByID, "ChangeBy", func() packets.Packet {
		return NewChangeByPacket()
	})
}
