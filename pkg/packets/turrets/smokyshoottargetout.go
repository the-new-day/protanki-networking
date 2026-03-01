package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Our smoky shot hit a target
type SmokyShootTargetOutPacket struct {
	packets.BasePacket
}

func NewSmokyShootTargetOutPacket() *SmokyShootTargetOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
	}

	attributes := []string{
		"clientTime",
		"target",
		"incarnationID",
		"targetBodyPosition",
		"localHitPoint",
		"globalHitPoint",
	}

	var id int32 = 229267683

	return &SmokyShootTargetOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(229267683, "SmokyShootTargetOut", func() packets.Packet {
		return NewSmokyShootTargetOutPacket()
	})
}
