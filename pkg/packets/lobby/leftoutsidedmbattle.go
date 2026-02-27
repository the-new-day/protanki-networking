package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A player has left a DM battle, the observer being outside the battle
type LeftOutsideDmBattlePacket struct {
	packets.BasePacket
}

func NewLeftOutsideDmBattlePacket() *LeftOutsideDmBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
		"username",
	}

	var id int32 = 504016996

	return &LeftOutsideDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(504016996, "LeftOutsideDmBattlePacket", func() packets.Packet {
		return NewLeftOutsideDmBattlePacket()
	})
}
