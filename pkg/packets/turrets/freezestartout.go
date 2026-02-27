package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// When we start using Freeze
type FreezeStartOutPacket struct {
	packets.BasePacket
}

func NewFreezeStartOutPacket() *FreezeStartOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -75406982

	return &FreezeStartOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-75406982, "FreezeStartOutPacket", func() packets.Packet {
		return NewFreezeStartOutPacket()
	})
}
