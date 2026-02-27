package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle user stats
type BattleUserTeamStatsPacket struct {
	packets.BasePacket
}

func NewBattleUserTeamStatsPacket() *BattleUserTeamStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(custom.NewBattleUserStatsCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"usersStat",
		"team",
	}

	var id int32 = -497293992

	return &BattleUserTeamStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-497293992, "BattleUserTeamStatsPacket", func() packets.Packet {
		return NewBattleUserTeamStatsPacket()
	})
}
