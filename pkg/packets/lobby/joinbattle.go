package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -1284211503

	return &JoinBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1284211503, "JoinBattlePacket", func() packets.Packet {
        return NewJoinBattlePacket()
    })
}
