package complex

import "github.com/the-new-day/probogo/internal/codec/factory"

type VectorStringCodec struct {
	factory.VectorCodec[string]
}

func NewVectorStringCodec() *VectorStringCodec {
	return &VectorStringCodec{
		VectorCodec: *factory.NewVectorCodec(&StringCodec{}, true),
	}
}
