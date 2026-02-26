package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
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

	var id int32 = 1478921140

	return &SmokyShootAirOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
