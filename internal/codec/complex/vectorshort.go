package complex

import (
	"github.com/the-new-day/probogo/internal/codec/factory"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type VectorShortCodec struct {
	factory.VectorCodec[int16]
}

func NewVectorShortCodec() *VectorShortCodec {
	return &VectorShortCodec{
		VectorCodec: *factory.NewVectorCodec(&primitive.ShortCodec{}, true),
	}
}
