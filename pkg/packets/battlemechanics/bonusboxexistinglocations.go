package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Locations of existing bonus boxes
type BonusBoxExistingLocationsPacket struct {
	packets.BasePacket
}

func NewBonusBoxExistingLocationsPacket() *BonusBoxExistingLocationsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.BonusBoxExistingLocationsID

	return &BonusBoxExistingLocationsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BonusBoxExistingLocationsID, "BonusBoxExistingLocations", func() packets.Packet {
		return NewBonusBoxExistingLocationsPacket()
	})
}
