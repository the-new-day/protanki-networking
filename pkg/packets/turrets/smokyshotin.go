package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.SmokyShotInID

	return &SmokyShotInPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SmokyShotInID, "SmokyShotIn", func() packets.Packet {
		return NewSmokyShotInPacket()
	})
}
