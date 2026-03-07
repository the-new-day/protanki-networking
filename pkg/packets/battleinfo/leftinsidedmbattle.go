package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.LeftInsideDmBattleID

	return &LeftInsideDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LeftInsideDmBattleID, "LeftInsideDmBattle", func() packets.Packet {
		return NewLeftInsideDmBattlePacket()
	})
}
