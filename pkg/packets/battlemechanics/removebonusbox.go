package battlemechanics

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Supply box bonus id
type RemoveBonusBoxPacket struct {
	packets.BasePacket
}

func NewRemoveBonusBoxPacket() *RemoveBonusBoxPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"bonusId",
	}

	var id int32 = -2026749922

	return &RemoveBonusBoxPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
