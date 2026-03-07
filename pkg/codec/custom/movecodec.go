package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type MoveCodec struct {
	CustomCodec
}

func NewMoveCodec() *MoveCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "angV", complex.NewVector3DCodec())
	AddField(customCodec, "control", &primitive.ByteCodec{})
	AddField(customCodec, "linV", complex.NewVector3DCodec())
	AddField(customCodec, "orientation", complex.NewVector3DCodec())
	AddField(customCodec, "pos", complex.NewVector3DCodec())

	return &MoveCodec{CustomCodec: *customCodec}
}
