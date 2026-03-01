package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Reason why player was kicked from battle
type BattleKickReasonPacket struct {
	packets.BasePacket
}

func NewBattleKickReasonPacket() *BattleKickReasonPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"reason",
	}

	var id int32 = -322235316

	return &BattleKickReasonPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-322235316, "BattleKickReason", func() packets.Packet {
		return NewBattleKickReasonPacket()
	})
}
