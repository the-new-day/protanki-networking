package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// The existing battle round has finished
type RoundFinishPacket struct {
	packets.BasePacket
}

func NewRoundFinishPacket() *RoundFinishPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = packets.RoundFinishID

	return &RoundFinishPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RoundFinishID, "RoundFinish", func() packets.Packet {
		return NewRoundFinishPacket()
	})
}
