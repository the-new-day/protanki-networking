package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -173682854

	return &RespawnDelayPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-173682854, "RespawnDelayPacket", func() packets.Packet {
        return NewRespawnDelayPacket()
    })
}
