package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
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

	var id int32 = packets.FireDamageID

	return &FireDamagePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FireDamageID, "FireDamage", func() packets.Packet {
		return NewFireDamagePacket()
	})
}
