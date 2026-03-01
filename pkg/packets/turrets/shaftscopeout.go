package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends server details about a released Shaft scope shot
type ShaftScopeOutPacket struct {
	packets.BasePacket
}

func NewShaftScopeOutPacket() *ShaftScopeOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
		codec.Wrap(complex.NewVectorShortCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
	}

	attributes := []string{
		"clientTime",
		"staticHitPoint",
		"targets",
		"targetHitPoints",
		"incarnationIDs",
		"targetBodyPositions",
		"globalHitPoints",
	}

	var id int32 = 1632423559

	return &ShaftScopeOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1632423559, "ShaftScopeOut", func() packets.Packet {
		return NewShaftScopeOutPacket()
	})
}
