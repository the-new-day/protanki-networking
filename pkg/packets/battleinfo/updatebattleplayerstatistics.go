package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates the in-battle statistics of a player.
type UpdateBattlePlayerStatisticsPacket struct {
	packets.BasePacket
}

func NewUpdateBattlePlayerStatisticsPacket() *UpdateBattlePlayerStatisticsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserStatsCodec(), false)),
	}

	attributes := []string{
		"userStats",
	}

	var id int32 = 696140460

	return &UpdateBattlePlayerStatisticsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
