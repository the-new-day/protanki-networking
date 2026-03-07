package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.BattleKickReasonID

	return &BattleKickReasonPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleKickReasonID, "BattleKickReason", func() packets.Packet {
		return NewBattleKickReasonPacket()
	})
}
