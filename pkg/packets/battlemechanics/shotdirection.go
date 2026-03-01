package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Shot direction
type ShotDirectionPacket struct {
	packets.BasePacket
}

func NewShotDirectionPacket() *ShotDirectionPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
	}

	attributes := []string{
		"shooter",
		"shotDirectionX",
		"shotDirectionY",
		"shotDirectionZ",
	}

	var id int32 = -118119523

	return &ShotDirectionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-118119523, "ShotDirection", func() packets.Packet {
		return NewShotDirectionPacket()
	})
}
