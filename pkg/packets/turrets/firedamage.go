package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Firebird attack
type FireDamagePacket struct {
	packets.BasePacket
}

func NewFireDamagePacket() *FireDamagePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorShortCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
	}

	attributes := []string{
		"clientTime",
		"targets",
		"incarnationIDs",
		"targetBodyPositions",
		"targetHitPoints",
	}

	var id int32 = 1395251766

	return &FireDamagePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
