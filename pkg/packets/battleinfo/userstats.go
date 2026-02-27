package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle User Stats
type BattleUserStatsPacket struct {
	packets.BasePacket
}

func NewBattleUserStatsPacket() *BattleUserStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserStatsCodec(), false)),
	}

	attributes := []string{
		"usersStat",
	}

	var id int32 = 1061006142

	return &BattleUserStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1061006142, "BattleUserStatsPacket", func() packets.Packet {
		return NewBattleUserStatsPacket()
	})
}
