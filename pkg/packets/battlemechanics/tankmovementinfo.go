package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client moved passively
type TankMovementInfoPacket struct {
	packets.BasePacket
}

func NewTankMovementInfoPacket() *TankMovementInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(custom.NewMoveCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
	}

	attributes := []string{
		"clientTime",
		"specificationID",
		"movement",
		"turretDirection",
	}

	var id int32 = -1683279062

	return &TankMovementInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1683279062, "TankMovementInfo", func() packets.Packet {
		return NewTankMovementInfoPacket()
	})
}
