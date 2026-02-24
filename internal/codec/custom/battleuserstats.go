package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type BattleUserStatsCodec struct {
	CustomCodec
}

func NewBattleUserStatsCodec() *BattleUserStatsCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "deaths", &primitive.ShortCodec{})
	AddField(customCodec, "kills", &primitive.ShortCodec{})
	AddField(customCodec, "score", &primitive.IntCodec{})
	AddField(customCodec, "username", &complex.StringCodec{})

	return &BattleUserStatsCodec{CustomCodec: *customCodec}
}
