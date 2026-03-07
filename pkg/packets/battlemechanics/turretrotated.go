package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Turret rotation data from server
type TurretRotatedPacket struct {
	packets.BasePacket
}

func NewTurretRotatedPacket() *TurretRotatedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewTurretRotateCodec()),
	}

	attributes := []string{
		"username",
		"turretRotation",
	}

	var id int32 = packets.TurretRotatedID

	return &TurretRotatedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TurretRotatedID, "TurretRotated", func() packets.Packet {
		return NewTurretRotatedPacket()
	})
}
