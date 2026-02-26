package complex

import (
	"github.com/the-new-day/probogo/pkg/codec/multiple"
)

type VectorStringCodec struct {
	multiple.VectorCodec[string]
}

func NewVectorStringCodec() *VectorStringCodec {
	return &VectorStringCodec{
		VectorCodec: *multiple.NewVectorCodec(NewStringCodec(), true),
	}
}
