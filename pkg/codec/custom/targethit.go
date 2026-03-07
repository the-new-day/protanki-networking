package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type TargetHitCodec struct {
	CustomCodec
}

func NewTargetHitCodec() *TargetHitCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "direction", complex.NewVector3DCodec())
	AddField(customCodec, "localHitPoint", complex.NewVector3DCodec())
	AddField(customCodec, "numberOfHits", &primitive.ByteCodec{})
	AddField(customCodec, "target", complex.NewStringCodec())

	return &TargetHitCodec{CustomCodec: *customCodec}
}
