package complex

import (
	"github.com/the-new-day/probogo/internal/codecs/multiple"
	"github.com/the-new-day/probogo/internal/codecs/primitive"
)

type VectorShortCodec struct {
	multiple.VectorCodec[int16]
}

func NewVectorShortCodec() *VectorShortCodec {
	return &VectorShortCodec{
		VectorCodec: *multiple.NewVectorCodec(&primitive.ShortCodec{}, true),
	}
}
