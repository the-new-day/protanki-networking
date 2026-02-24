package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle user rewards.
type UserRewardsPacket struct {
	packets.BasePacket
}

func NewUserRewardsPacket() *UserRewardsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserRewardsCodec(), false)),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{
		"reward",
		"timeToRestart",
	}

	var id int32 = 560336625

	return &UserRewardsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
