package complex

import "github.com/the-new-day/probogo/internal/codecs/multiple"

type VectorVector3DCodec struct {
	multiple.VectorCodec[map[string]any]
}

func NewVectorVector3DCodec() *VectorVector3DCodec {
	return &VectorVector3DCodec{
		VectorCodec: *multiple.NewVectorCodec(NewVector3DCodec(), true),
	}
}
