package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
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

	var id int32 = packets.BattleFundID

	return &BattleFundPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleFundID, "BattleFund", func() packets.Packet {
		return NewBattleFundPacket()
	})
}
