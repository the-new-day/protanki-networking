package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
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

	var id int32 = -1749108178

	return &TurretControlPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
