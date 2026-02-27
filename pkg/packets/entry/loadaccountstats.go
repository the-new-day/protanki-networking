package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads the player
type LoadAccountStatsPacket struct {
	packets.BasePacket
}

func NewLoadAccountStatsPacket() *LoadAccountStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.ByteCodec{}),
		codec.Wrap(&primitive.FloatCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"crystals",
		"currentRankXP",
		"doubleCrystalLength",
		"hasDoubleCrystal",
		"nextRankXP",
		"ratingRank",
		"rank",
		"rating",
		"XP",
		"serverNumber",
		"username",
		"userProfileUrl",
	}

	var id int32 = 907073245

	return &LoadAccountStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(907073245, "LoadAccountStatsPacket", func() packets.Packet {
		return NewLoadAccountStatsPacket()
	})
}
