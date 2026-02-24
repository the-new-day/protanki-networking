package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Update the health of a tank.
type TankHealthPacket struct {
	packets.BasePacket
}

func NewTankHealthPacket() *TankHealthPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{
		"username",
		"health",
	}

	var id int32 = -611961116

	return &TankHealthPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
