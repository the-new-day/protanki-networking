package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Server syncs tank movement stats with client
type TankStatSyncedPacket struct {
	packets.BasePacket
}

func NewTankStatSyncedPacket() *TankStatSyncedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
	}

	attributes := []string{
		"username",
		"maxSpeed",
		"maxTurnSpeed",
		"maxTurretRotationSpeed",
		"acceleration",
		"specificationID",
	}

	var id int32 = packets.TankStatSyncedID

	return &TankStatSyncedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TankStatSyncedID, "TankStatSynced", func() packets.Packet {
		return NewTankStatSyncedPacket()
	})
}
