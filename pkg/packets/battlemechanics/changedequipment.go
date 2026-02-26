package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// User Changed Equipment
type UserChangedEquipmentPacket struct {
	packets.BasePacket
}

func NewUserChangedEquipmentPacket() *UserChangedEquipmentPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"tank",
	}

	var id int32 = -1767633906

	return &UserChangedEquipmentPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
