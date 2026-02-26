package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A tank has been killed
type KillConfirmPacket struct {
	packets.BasePacket
}

func NewKillConfirmPacket() *KillConfirmPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"target",
		"killer",
		"respDelay",
	}

	var id int32 = -42520728

	return &KillConfirmPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
