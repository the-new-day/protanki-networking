package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A player has left a DM battle, the observer being within the battle
type LeftInsideDmBattlePacket struct {
	packets.BasePacket
}

func NewLeftInsideDmBattlePacket() *LeftInsideDmBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = -1689876764

	return &LeftInsideDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1689876764, "LeftInsideDmBattlePacket", func() packets.Packet {
        return NewLeftInsideDmBattlePacket()
    })
}
