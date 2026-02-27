package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Creates a new battle
type CreateBattlePacket struct {
	packets.BasePacket
}

func NewCreateBattlePacket() *CreateBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(custom.NewBattleLimitsCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(custom.NewRankRangeCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"autoBalance",
		"battleMode",
		"format",
		"friendlyFire",
		"battleLimits",
		"mapID",
		"maxPeopleCount",
		"name",
		"parkourMode",
		"privateBattle",
		"proBattle",
		"rankRange",
		"noRearm",
		"theme",
		"noSupplyBoxes",
		"noCrystalBoxes",
		"noSupplies",
		"noUpgrade",
		"lowResistance",
		"useDropTimings",
		"noGoldBoxes",
		"noGoldSiren",
		"noGoldDropZone",
		"noMedkit",
		"noMines",
		"randomGoldChance",
		"dependentCooldown",
	}

	var id int32 = -2135234426

	return &CreateBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-2135234426, "CreateBattlePacket", func() packets.Packet {
        return NewCreateBattlePacket()
    })
}
