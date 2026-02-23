package complex

import (
	"github.com/the-new-day/probogo/internal/codec/factory"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type Vector3DCodec struct {
	factory.MultiCodec[float32]
}

func NewVector3DCodec() *Vector3DCodec {
	return &Vector3DCodec{
		MultiCodec: *factory.NewMultiCodec([]string{"x", "y", "z"}, &primitive.FloatCodec{}, true),
	}
}
