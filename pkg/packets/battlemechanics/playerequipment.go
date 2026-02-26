package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player Equipment.
type PlayerEquipmentPacket struct {
	packets.BasePacket
}

func NewPlayerEquipmentPacket() *PlayerEquipmentPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -1643824092

	return &PlayerEquipmentPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
