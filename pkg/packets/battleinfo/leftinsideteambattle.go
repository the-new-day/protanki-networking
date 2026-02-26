package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A player has left a team battle, the observer being within the battle
type LeftInsideTeamBattlePacket struct {
	packets.BasePacket
}

func NewLeftInsideTeamBattlePacket() *LeftInsideTeamBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1411656080

	return &LeftInsideTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
