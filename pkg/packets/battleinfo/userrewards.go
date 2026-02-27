package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle User Rewards
type BattleUserRewardsPacket struct {
	packets.BasePacket
}

func NewBattleUserRewardsPacket() *BattleUserRewardsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserRewardsCodec(), false)),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"reward",
		"timeToRestart",
	}

	var id int32 = 560336625

	return &BattleUserRewardsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(560336625, "BattleUserRewardsPacket", func() packets.Packet {
		return NewBattleUserRewardsPacket()
	})
}
