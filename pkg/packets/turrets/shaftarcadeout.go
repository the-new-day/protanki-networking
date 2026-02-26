package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends server details about a released Shaft arcade shot
type ShaftArcadeOutPacket struct {
	packets.BasePacket
}

func NewShaftArcadeOutPacket() *ShaftArcadeOutPacket {
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

	var id int32 = -2030760866

	return &ShaftArcadeOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
