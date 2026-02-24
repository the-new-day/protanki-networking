package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
)

type BattleNotifierCodec struct {
	CustomCodec
}

func NewBattleNotifierCodec() *BattleNotifierCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "battleInfo", &BattleInfoCodec{})
	AddField(customCodec, "username", &complex.StringCodec{})

	return &BattleNotifierCodec{CustomCodec: *customCodec}
}
