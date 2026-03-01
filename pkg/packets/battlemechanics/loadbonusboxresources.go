package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Bonus Box Resources
type LoadBonusBoxResourcesPacket struct {
	packets.BasePacket
}

func NewLoadBonusBoxResourcesPacket() *LoadBonusBoxResourcesPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.LoadBonusBoxResourcesID

	return &LoadBonusBoxResourcesPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadBonusBoxResourcesID, "LoadBonusBoxResources", func() packets.Packet {
		return NewLoadBonusBoxResourcesPacket()
	})
}
