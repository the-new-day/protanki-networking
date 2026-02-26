package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player Rating stats
type LoadRatingStatsPacket struct {
	packets.BasePacket
}

func NewLoadRatingStatsPacket() *LoadRatingStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"rating",
		"place",
	}

	var id int32 = -1128606444

	return &LoadRatingStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
