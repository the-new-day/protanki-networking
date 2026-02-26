package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Flag has been delivered
type FlagDeliveredPacket struct {
	packets.BasePacket
}

func NewFlagDeliveredPacket() *FlagDeliveredPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"baseteam",
		"username",
	}

	var id int32 = -1870108387

	return &FlagDeliveredPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
