package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.FreezeEndOutID

	return &FreezeEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FreezeEndOutID, "FreezeEndOut", func() packets.Packet {
		return NewFreezeEndOutPacket()
	})
}
