package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Damage dealt to a tank.
type TankDamagePacket struct {
	packets.BasePacket
}

func NewTankDamagePacket() *TankDamagePacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewTankDamageCodec(), false)),
	}
	attributes := []string{
		"damages",
	}

	var id int32 = -1165230470

	return &TankDamagePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
