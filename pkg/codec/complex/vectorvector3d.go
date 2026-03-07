package complex

import "github.com/the-new-day/protanki-networking/pkg/codec/multiple"

type VectorVector3DCodec struct {
	multiple.VectorCodec[map[string]float32]
}

func NewVectorVector3DCodec() *VectorVector3DCodec {
	return &VectorVector3DCodec{
		VectorCodec: *multiple.NewVectorCodec(NewVector3DCodec(), true),
	}
}
