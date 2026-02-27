package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends current turret rotation data to the server
type TurretRotationPacket struct {
	packets.BasePacket
}

func NewTurretRotationPacket() *TurretRotationPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(custom.NewTurretRotateCodec()),
		codec.Wrap(&primitive.ShortCodec{}),
	}

	attributes := []string{
		"clientTime",
		"turretRotation",
		"incarnationID",
	}

	var id int32 = -114968993

	return &TurretRotationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-114968993, "TurretRotationPacket", func() packets.Packet {
        return NewTurretRotationPacket()
    })
}
