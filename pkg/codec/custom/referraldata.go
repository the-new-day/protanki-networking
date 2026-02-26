package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type ReferralDataCodec struct {
	CustomCodec
}

func NewReferralDataCodec() *ReferralDataCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "income", &primitive.IntCodec{})
	AddField(customCodec, "username", complex.NewStringCodec())

	return &ReferralDataCodec{CustomCodec: *customCodec}
}
