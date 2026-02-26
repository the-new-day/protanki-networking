package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player Ranked Up
type RankUpPacket struct {
	packets.BasePacket
}

func NewRankUpPacket() *RankUpPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"rank",
	}

	var id int32 = 1262947513

	return &RankUpPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
