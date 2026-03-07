package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Client requests to join the selected battle
type JoinBattlePacket struct {
	packets.BasePacket
}

func NewJoinBattlePacket() *JoinBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"team",
	}

	var id int32 = packets.JoinBattleID

	return &JoinBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.JoinBattleID, "JoinBattle", func() packets.Packet {
		return NewJoinBattlePacket()
	})
}
