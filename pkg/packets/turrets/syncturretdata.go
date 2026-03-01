package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Syncs turret data to the client
type SyncTurretDataPacket struct {
	packets.BasePacket
}

func NewSyncTurretDataPacket() *SyncTurretDataPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -2124388778

	return &SyncTurretDataPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-2124388778, "SyncTurretData", func() packets.Packet {
		return NewSyncTurretDataPacket()
	})
}
