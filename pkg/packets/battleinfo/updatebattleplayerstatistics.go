package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Update the in-battle statistics of a player.
type UpdateBattlePlayerStatisticsPacket struct {
	packets.BasePacket
}

func NewUpdateBattlePlayerStatisticsPacket() *UpdateBattlePlayerStatisticsPacket {
	codecs := []codec.Codec{
		codec.Wrap(custom.NewBattleUserStatsCodec()),
	}
	attributes := []string{
		"userStats",
	}

	var id int32 = 696140460

	return &UpdateBattlePlayerStatisticsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
