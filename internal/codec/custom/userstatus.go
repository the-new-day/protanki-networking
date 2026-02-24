package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type UserStatusCodec struct {
	CustomCodec
}

func NewUserStatusCodec() *UserStatusCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "modLevel", &primitive.IntCodec{})
	AddField(customCodec, "ip", complex.NewStringCodec())
	AddField(customCodec, "rank", &primitive.IntCodec{})
	AddField(customCodec, "username", complex.NewStringCodec())

	return &UserStatusCodec{CustomCodec: *customCodec}
}
