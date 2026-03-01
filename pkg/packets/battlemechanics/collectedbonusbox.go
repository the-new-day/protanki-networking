package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A bonus box was picked up
type CollectedBonusBoxPacket struct {
	packets.BasePacket
}

func NewCollectedBonusBoxPacket() *CollectedBonusBoxPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"bonusId",
	}

	var id int32 = -1291499147

	return &CollectedBonusBoxPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1291499147, "CollectedBonusBox", func() packets.Packet {
		return NewCollectedBonusBoxPacket()
	})
}
