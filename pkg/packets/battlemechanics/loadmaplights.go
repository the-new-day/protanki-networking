package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Map Lights
type LoadMapLightsPacket struct {
	packets.BasePacket
}

func NewLoadMapLightsPacket() *LoadMapLightsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -152638117

	return &LoadMapLightsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-152638117, "LoadMapLightsPacket", func() packets.Packet {
		return NewLoadMapLightsPacket()
	})
}
