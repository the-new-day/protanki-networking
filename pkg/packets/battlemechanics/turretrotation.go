package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.TurretRotationID

	return &TurretRotationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TurretRotationID, "TurretRotation", func() packets.Packet {
		return NewTurretRotationPacket()
	})
}
