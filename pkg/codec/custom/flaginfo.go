package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
)

type FlagInfoCodec struct {
	CustomCodec
}

func NewFlagInfoCodec() *FlagInfoCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "pole_pos", complex.NewVector3DCodec())
	AddField(customCodec, "holder", complex.NewStringCodec())
	AddField(customCodec, "current_pos", complex.NewVector3DCodec())

	return &FlagInfoCodec{CustomCodec: *customCodec}
}
