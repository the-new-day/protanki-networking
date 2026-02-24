package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A player has left a DM battle, the observer being within the battle.
type LeftInsideDMBattlePacket struct {
	packets.BasePacket
}

func NewLeftInsideDMBattlePacket() *LeftInsideDMBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}
	attributes := []string{
		"username",
	}

	var id int32 = -1689876764

	return &LeftInsideDMBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
