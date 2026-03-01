package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Packet sent when the player has fully respawned
type FullyRespawnedPacket struct {
	packets.BasePacket
}

func NewFullyRespawnedPacket() *FullyRespawnedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = packets.FullyRespawnedID

	return &FullyRespawnedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FullyRespawnedID, "FullyRespawned", func() packets.Packet {
		return NewFullyRespawnedPacket()
	})
}
