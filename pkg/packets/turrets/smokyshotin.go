package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Smokey Shoot in
type SmokyShotInPacket struct {
	packets.BasePacket
}

func NewSmokyShotInPacket() *SmokyShotInPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"shooter",
		"target",
		"hitPoint",
		"weakeningCoeff",
		"isCritical",
	}

	var id int32 = -1334002026

	return &SmokyShotInPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
