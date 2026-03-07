package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Loads current battle statistics
type InitBattleStatsPacket struct {
	packets.BasePacket
}

func NewInitBattleStatsPacket() *InitBattleStatsPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(custom.NewBattleLimitsCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleMode",
		"format",
		"fund",
		"battleLimits",
		"mapName",
		"maxPeopleCount",
		"parkourMode",
		"premiumBonusPercentage",
		"spectator",
		"suspiciousUsers",
		"timeLeft",
	}

	var id int32 = packets.InitBattleStatsID

	return &InitBattleStatsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.InitBattleStatsID, "InitBattleStats", func() packets.Packet {
		return NewInitBattleStatsPacket()
	})
}
