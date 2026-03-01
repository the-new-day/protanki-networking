package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates the health of a tank
type TankHealthPacket struct {
	packets.BasePacket
}

func NewTankHealthPacket() *TankHealthPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
	}

	attributes := []string{
		"username",
		"health",
	}

	var id int32 = packets.TankHealthID

	return &TankHealthPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TankHealthID, "TankHealth", func() packets.Packet {
		return NewTankHealthPacket()
	})
}
