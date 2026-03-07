package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Battle user control packet
type BattleUserControlPacket struct {
	packets.BasePacket
}

func NewBattleUserControlPacket() *BattleUserControlPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.ByteCodec{}),
	}

	attributes := []string{
		"tankiId",
		"control",
	}

	var id int32 = packets.BattleUserControlID

	return &BattleUserControlPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleUserControlID, "BattleUserControl", func() packets.Packet {
		return NewBattleUserControlPacket()
	})
}
