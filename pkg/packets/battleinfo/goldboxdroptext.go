package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Announces that a gold box will drop soon
type GoldBoxDropTextPacket struct {
	packets.BasePacket
}

func NewGoldBoxDropTextPacket() *GoldBoxDropTextPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"text",
		"soundID",
	}

	var id int32 = -666893269

	return &GoldBoxDropTextPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-666893269, "GoldBoxDropTextPacket", func() packets.Packet {
		return NewGoldBoxDropTextPacket()
	})
}
