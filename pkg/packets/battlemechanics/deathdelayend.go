package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Respawn delay ends and player can start respawning
type DeathDelayEndPacket struct {
	packets.BasePacket
}

func NewDeathDelayEndPacket() *DeathDelayEndPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.DeathDelayEndID

	return &DeathDelayEndPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.DeathDelayEndID, "DeathDelayEnd", func() packets.Packet {
		return NewDeathDelayEndPacket()
	})
}
