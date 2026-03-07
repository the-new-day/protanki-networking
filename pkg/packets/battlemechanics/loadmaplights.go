package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.LoadMapLightsID

	return &LoadMapLightsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadMapLightsID, "LoadMapLights", func() packets.Packet {
		return NewLoadMapLightsPacket()
	})
}
