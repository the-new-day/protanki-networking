package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type UserStatusCodec struct {
	CustomCodec
}

func NewUserStatusCodec() *UserStatusCodec {
	customCodec := NewCustomCodec(true)

	AddField(customCodec, "modLevel", &primitive.IntCodec{})
	AddField(customCodec, "ip", complex.NewStringCodec())
	AddField(customCodec, "rank", &primitive.IntCodec{})
	AddField(customCodec, "username", complex.NewStringCodec())

	return &UserStatusCodec{CustomCodec: *customCodec}
}
