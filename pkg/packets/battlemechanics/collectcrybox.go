package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Collect a crystal box
type CollectCryBoxPacket struct {
	packets.BasePacket
}

func NewCollectCryBoxPacket() *CollectCryBoxPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"bonusId",
	}

	var id int32 = -1047185003

	return &CollectCryBoxPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1047185003, "CollectCryBoxPacket", func() packets.Packet {
        return NewCollectCryBoxPacket()
    })
}
