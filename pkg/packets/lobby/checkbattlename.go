package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Censors invalid battle names
type CheckBattleNamePacket struct {
	packets.BasePacket
}

func NewCheckBattleNamePacket() *CheckBattleNamePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleName",
	}

	var id int32 = 566652736

	return &CheckBattleNamePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
