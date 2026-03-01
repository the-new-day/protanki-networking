package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Turret Control Packet
type TurretControlPacket struct {
	packets.BasePacket
}

func NewTurretControlPacket() *TurretControlPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(&primitive.ByteCodec{}),
	}

	attributes := []string{
		"clientTime",
		"specificationID",
		"control",
	}

	var id int32 = packets.TurretControlID

	return &TurretControlPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TurretControlID, "TurretControl", func() packets.Packet {
		return NewTurretControlPacket()
	})
}
