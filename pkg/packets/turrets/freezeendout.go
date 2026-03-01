package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// When we stop using Freeze
type FreezeEndOutPacket struct {
	packets.BasePacket
}

func NewFreezeEndOutPacket() *FreezeEndOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -1654947652

	return &FreezeEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1654947652, "FreezeEndOut", func() packets.Packet {
		return NewFreezeEndOutPacket()
	})
}
