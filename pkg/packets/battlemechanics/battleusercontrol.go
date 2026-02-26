package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -301298508

	return &BattleUserControlPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
