package status

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -962759489

	return &RankStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-962759489, "RankStatusPacket", func() packets.Packet {
        return NewRankStatusPacket()
    })
}
