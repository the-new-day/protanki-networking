package complex

import (
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type VectorShortCodec struct {
	multiple.VectorCodec[int16]
}

func NewVectorShortCodec() *VectorShortCodec {
	return &VectorShortCodec{
		VectorCodec: *multiple.NewVectorCodec(&primitive.ShortCodec{}, true),
	}
}
