package custom

import (
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type TurretRotateCodec struct {
	CustomCodec
}

func NewTurretRotateCodec() *TurretRotateCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "angle", &primitive.FloatCodec{})
	AddField(customCodec, "control", &primitive.ByteCodec{})

	return &TurretRotateCodec{CustomCodec: *customCodec}
}
