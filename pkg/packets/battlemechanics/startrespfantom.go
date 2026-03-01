package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Information about fantom status of a player
type StartRespFantomPacket struct {
	packets.BasePacket
}

func NewStartRespFantomPacket() *StartRespFantomPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(&primitive.ShortCodec{}),
		codec.Wrap(&primitive.ShortCodec{}),
	}

	attributes := []string{
		"username",
		"team",
		"position",
		"orientation",
		"health",
		"incarnationID",
	}

	var id int32 = packets.StartRespFantomID

	return &StartRespFantomPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.StartRespFantomID, "StartRespFantom", func() packets.Packet {
		return NewStartRespFantomPacket()
	})
}
