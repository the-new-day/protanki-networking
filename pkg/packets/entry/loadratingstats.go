package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.LoadRatingStatsID

	return &LoadRatingStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadRatingStatsID, "LoadRatingStats", func() packets.Packet {
		return NewLoadRatingStatsPacket()
	})
}
