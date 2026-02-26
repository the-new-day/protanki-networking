package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type BattleUserRewardsCodec struct {
	CustomCodec
}

func NewBattleUserRewardsCodec() *BattleUserRewardsCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "newbiesAbonementBonusReward", &primitive.IntCodec{})
	AddField(customCodec, "premiumBonusReward", &primitive.IntCodec{})
	AddField(customCodec, "reward", &primitive.IntCodec{})
	AddField(customCodec, "userid", complex.NewStringCodec())

	return &BattleUserRewardsCodec{CustomCodec: *customCodec}
}
