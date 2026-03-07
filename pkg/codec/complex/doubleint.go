package complex

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type DoubleIntCodec struct {
	multiple.MultiCodec[int32]
}

func NewDoubleIntCodec(attr1 string, attr2 string) *DoubleIntCodec {
	return &DoubleIntCodec{
		MultiCodec: *multiple.NewMultiCodec([]string{attr1, attr2}, &primitive.IntCodec{}, false),
	}
}
