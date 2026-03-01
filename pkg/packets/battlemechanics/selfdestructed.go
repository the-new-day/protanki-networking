package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Information about a self destructed tank
type SelfDestructedPacket struct {
	packets.BasePacket
}

func NewSelfDestructedPacket() *SelfDestructedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"respDelay",
	}

	var id int32 = packets.SelfDestructedID

	return &SelfDestructedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SelfDestructedID, "SelfDestructed", func() packets.Packet {
		return NewSelfDestructedPacket()
	})
}
