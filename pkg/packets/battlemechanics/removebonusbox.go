package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
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

	var id int32 = packets.RemoveBonusBoxID

	return &RemoveBonusBoxPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RemoveBonusBoxID, "RemoveBonusBox", func() packets.Packet {
		return NewRemoveBonusBoxPacket()
	})
}
