package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Respawn delay ends and player can start respawning
type DeathDelayEndPacket struct {
	packets.BasePacket
}

func NewDeathDelayEndPacket() *DeathDelayEndPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 268832557

	return &DeathDelayEndPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
