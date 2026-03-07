package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Firebird stops shooting
type FireEndOutPacket struct {
	packets.BasePacket
}

func NewFireEndOutPacket() *FireEndOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = packets.FireEndOutID

	return &FireEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FireEndOutID, "FireEndOut", func() packets.Packet {
		return NewFireEndOutPacket()
	})
}
