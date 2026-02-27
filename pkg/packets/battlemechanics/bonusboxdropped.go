package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A bonus box has dropped
type BonusBoxDroppedPacket struct {
	packets.BasePacket
}

func NewBonusBoxDroppedPacket() *BonusBoxDroppedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"bonusId",
		"position",
		"fallTimeThreshold",
	}

	var id int32 = 1831462385

	return &BonusBoxDroppedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1831462385, "BonusBoxDroppedPacket", func() packets.Packet {
		return NewBonusBoxDroppedPacket()
	})
}
