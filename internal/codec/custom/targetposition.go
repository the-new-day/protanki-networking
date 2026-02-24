package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type TargetPositionCodec struct {
	CustomCodec
}

func NewTargetPositionCodec() *TargetPositionCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "localHitPoint", complex.NewVector3DCodec())
	AddField(customCodec, "orientation", complex.NewVector3DCodec())
	AddField(customCodec, "position", complex.NewVector3DCodec())
	AddField(customCodec, "target", complex.NewStringCodec())
	AddField(customCodec, "turretAngle", &primitive.FloatCodec{})

	return &TargetPositionCodec{CustomCodec: *customCodec}
}
