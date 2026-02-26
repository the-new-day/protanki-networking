package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
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
