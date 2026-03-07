package complex

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type VectorShortCodec struct {
	multiple.VectorCodec[int16]
}

func NewVectorShortCodec() *VectorShortCodec {
	return &VectorShortCodec{
		VectorCodec: *multiple.NewVectorCodec(&primitive.ShortCodec{}, true),
	}
}
