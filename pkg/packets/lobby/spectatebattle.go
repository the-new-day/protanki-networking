package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client requests to spectate the selected battle
type SpectateBattlePacket struct {
	packets.BasePacket
}

func NewSpectateBattlePacket() *SpectateBattlePacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.SpectateBattleID

	return &SpectateBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SpectateBattleID, "SpectateBattle", func() packets.Packet {
		return NewSpectateBattlePacket()
	})
}
