package status

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Loads the rank of a player
type RankStatusPacket struct {
	packets.BasePacket
}

func NewRankStatusPacket() *RankStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.ByteCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"rank",
		"username",
	}

	var id int32 = packets.RankStatusID

	return &RankStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RankStatusID, "RankStatus", func() packets.Packet {
		return NewRankStatusPacket()
	})
}
