package complex

import (
	"github.com/the-new-day/probogo/internal/codec/factory"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type DoubleIntCodec struct {
	factory.MultiCodec[int32]
}

func NewDoubleIntCodec(attr1 string, attr2 string) *DoubleIntCodec {
	return &DoubleIntCodec{
		MultiCodec: *factory.NewMultiCodec([]string{attr1, attr2}, &primitive.IntCodec{}, true),
	}
}
