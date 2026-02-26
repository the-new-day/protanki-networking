package status

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
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
