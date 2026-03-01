package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
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

	var id int32 = packets.UserChangedEquipmentID

	return &UserChangedEquipmentPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UserChangedEquipmentID, "UserChangedEquipment", func() packets.Packet {
		return NewUserChangedEquipmentPacket()
	})
}
