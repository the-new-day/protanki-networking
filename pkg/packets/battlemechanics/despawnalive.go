package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player despawns while alive
type DespawnAlivePacket struct {
	packets.BasePacket
}

func NewDespawnAlivePacket() *DespawnAlivePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1719707347

	return &DespawnAlivePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1719707347, "DespawnAlive", func() packets.Packet {
		return NewDespawnAlivePacket()
	})
}
