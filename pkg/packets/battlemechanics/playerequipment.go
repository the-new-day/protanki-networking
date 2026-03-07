package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.PlayerEquipmentID

	return &PlayerEquipmentPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PlayerEquipmentID, "PlayerEquipment", func() packets.Packet {
		return NewPlayerEquipmentPacket()
	})
}
