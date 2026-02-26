package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle fund
type BattleFundPacket struct {
	packets.BasePacket
}

func NewBattleFundPacket() *BattleFundPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"fund",
	}

	var id int32 = 1149211509

	return &BattleFundPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
