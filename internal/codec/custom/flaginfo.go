package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
)

type FlagInfoCodec struct {
	CustomCodec
}

func NewFlagInfoCodec() *FlagInfoCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "pole_pos", &complex.Vector3DCodec{})
	AddField(customCodec, "holder", &complex.StringCodec{})
	AddField(customCodec, "current_pos", &complex.Vector3DCodec{})

	return &FlagInfoCodec{CustomCodec: *customCodec}
}
