package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Smokey Shoot at Air
type SmokyShootAirOutPacket struct {
	packets.BasePacket
}

func NewSmokyShootAirOutPacket() *SmokyShootAirOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = packets.SmokyShootAirOutID

	return &SmokyShootAirOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SmokyShootAirOutID, "SmokyShootAirOut", func() packets.Packet {
		return NewSmokyShootAirOutPacket()
	})
}
