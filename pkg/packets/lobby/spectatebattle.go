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

	var id int32 = -1315002220

	return &SpectateBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1315002220, "SpectateBattlePacket", func() packets.Packet {
        return NewSpectateBattlePacket()
    })
}
