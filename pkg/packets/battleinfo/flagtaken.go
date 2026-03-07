package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Flag has been taken
type FlagTakenPacket struct {
	packets.BasePacket
}

func NewFlagTakenPacket() *FlagTakenPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"flagteam",
	}

	var id int32 = packets.FlagTakenID

	return &FlagTakenPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.FlagTakenID, "FlagTaken", func() packets.Packet {
		return NewFlagTakenPacket()
	})
}
