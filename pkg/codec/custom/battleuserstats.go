package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type BattleUserStatsCodec struct {
	CustomCodec
}

func NewBattleUserStatsCodec() *BattleUserStatsCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "deaths", &primitive.ShortCodec{})
	AddField(customCodec, "kills", &primitive.ShortCodec{})
	AddField(customCodec, "score", &primitive.IntCodec{})
	AddField(customCodec, "username", complex.NewStringCodec())

	return &BattleUserStatsCodec{CustomCodec: *customCodec}
}
