package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type TankDamageCodec struct {
	CustomCodec
}

func NewTankDamageCodec() *TankDamageCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "damage", &primitive.IntCodec{})
	AddField(customCodec, "damageType", &primitive.IntCodec{})
	AddField(customCodec, "target", complex.NewStringCodec())

	return &TankDamageCodec{CustomCodec: *customCodec}
}
