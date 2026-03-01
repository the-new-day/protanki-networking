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

	var id int32 = packets.BattleUserRewardsID

	return &BattleUserRewardsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleUserRewardsID, "BattleUserRewards", func() packets.Packet {
		return NewBattleUserRewardsPacket()
	})
}
