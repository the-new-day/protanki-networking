package complex

import (
	"github.com/the-new-day/probogo/internal/codecs/multiple"
)

type VectorStringCodec struct {
	multiple.VectorCodec[string]
}

func NewVectorStringCodec() *VectorStringCodec {
	return &VectorStringCodec{
		VectorCodec: *multiple.NewVectorCodec(&StringCodec{}, true),
	}
}
