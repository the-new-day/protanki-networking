package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates the in-battle statistics of a player.
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

	var id int32 = packets.UpdateBattlePlayerStatisticsID

	return &UpdateBattlePlayerStatisticsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UpdateBattlePlayerStatisticsID, "UpdateBattlePlayerStatistics", func() packets.Packet {
		return NewUpdateBattlePlayerStatisticsPacket()
	})
}
