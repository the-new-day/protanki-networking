package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Starts a new round in the existing battle
type RoundStartPacket struct {
	packets.BasePacket
}

func NewRoundStartPacket() *RoundStartPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = packets.RoundStartID

	return &RoundStartPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RoundStartID, "RoundStart", func() packets.Packet {
		return NewRoundStartPacket()
	})
}
