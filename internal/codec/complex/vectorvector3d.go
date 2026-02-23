package complex

import "github.com/the-new-day/probogo/internal/codec/factory"

type VectorVector3DCodec struct {
	factory.VectorCodec[map[string]any]
}

func NewVectorVector3DCodec() *VectorVector3DCodec {
	return &VectorVector3DCodec{
		VectorCodec: *factory.NewVectorCodec(NewVector3DCodec(), true),
	}
}
