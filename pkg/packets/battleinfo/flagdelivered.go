package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.FlagDeliveredID

	return &FlagDeliveredPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FlagDeliveredID, "FlagDelivered", func() packets.Packet {
		return NewFlagDeliveredPacket()
	})
}
