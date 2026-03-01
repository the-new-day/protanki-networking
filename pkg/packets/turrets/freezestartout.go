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

	var id int32 = packets.FreezeStartOutID

	return &FreezeStartOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FreezeStartOutID, "FreezeStartOut", func() packets.Packet {
		return NewFreezeStartOutPacket()
	})
}
