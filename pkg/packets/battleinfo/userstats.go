package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle users stats.
type UserStatsPacket struct {
	packets.BasePacket
}

func NewUserStatsPacket() *UserStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserStatsCodec(), false)),
	}
	attributes := []string{
		"usersStats",
	}

	var id int32 = 1061006142

	return &UserStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
