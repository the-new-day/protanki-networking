package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
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

	var id int32 = packets.LeftInsideTeamBattleID

	return &LeftInsideTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LeftInsideTeamBattleID, "LeftInsideTeamBattle", func() packets.Packet {
		return NewLeftInsideTeamBattlePacket()
	})
}
