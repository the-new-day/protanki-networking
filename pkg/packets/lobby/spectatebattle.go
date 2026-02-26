package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client requests to spectate the selected battle
type SpectateBattlePacket struct {
	packets.BasePacket
}

func NewSpectateBattlePacket() *SpectateBattlePacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -1315002220

	return &SpectateBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
