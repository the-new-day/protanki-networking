package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Respawn Delay Packet
type RespawnDelayPacket struct {
	packets.BasePacket
}

func NewRespawnDelayPacket() *RespawnDelayPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"tank",
		"respawnDelay",
	}

	var id int32 = packets.RespawnDelayID

	return &RespawnDelayPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RespawnDelayID, "RespawnDelay", func() packets.Packet {
		return NewRespawnDelayPacket()
	})
}
