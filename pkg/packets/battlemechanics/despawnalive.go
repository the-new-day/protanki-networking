package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.DespawnAliveID

	return &DespawnAlivePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.DespawnAliveID, "DespawnAlive", func() packets.Packet {
		return NewDespawnAlivePacket()
	})
}
