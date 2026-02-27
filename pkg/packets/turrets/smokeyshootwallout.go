package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Smokey Shoot out
type SmokyShootWallOutPacket struct {
	packets.BasePacket
}

func NewSmokyShootWallOutPacket() *SmokyShootWallOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
	}

	attributes := []string{
		"clientTime",
		"hitPoint",
	}

	var id int32 = 1470597926

	return &SmokyShootWallOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1470597926, "SmokyShootWallOutPacket", func() packets.Packet {
        return NewSmokyShootWallOutPacket()
    })
}
