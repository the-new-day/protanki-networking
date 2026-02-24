package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type TargetHitCodec struct {
	CustomCodec
}

func NewTargetHitCodec() *TargetHitCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "direction", complex.NewVector3DCodec())
	AddField(customCodec, "localHitPoint", complex.NewVector3DCodec())
	AddField(customCodec, "numberOfHits", &primitive.ByteCodec{})
	AddField(customCodec, "target", &complex.StringCodec{})

	return &TargetHitCodec{CustomCodec: *customCodec}
}
