package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 1927704181

	return &TurretRotatedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1927704181, "TurretRotatedPacket", func() packets.Packet {
        return NewTurretRotatedPacket()
    })
}
