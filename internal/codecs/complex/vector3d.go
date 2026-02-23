package complex

import (
	"github.com/the-new-day/probogo/internal/codecs/multiple"
	"github.com/the-new-day/probogo/internal/codecs/primitive"
)

type Vector3DCodec struct {
	multiple.MultiCodec[float32]
}

func NewVector3DCodec() *Vector3DCodec {
	return &Vector3DCodec{
		MultiCodec: *multiple.NewMultiCodec([]string{"x", "y", "z"}, &primitive.FloatCodec{}, true),
	}
}
