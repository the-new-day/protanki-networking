package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sends server details about a railgun shot that has just started to release
type RailgunShotInitOutPacket struct {
	packets.BasePacket
}

func NewRailgunShotInitOutPacket() *RailgunShotInitOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = packets.RailgunShotInitOutID

	return &RailgunShotInitOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RailgunShotInitOutID, "RailgunShotInitOut", func() packets.Packet {
		return NewRailgunShotInitOutPacket()
	})
}
