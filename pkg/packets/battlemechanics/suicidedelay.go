package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Suicide delay packet
type SuicideDelayPacket struct {
	packets.BasePacket
}

func NewSuicideDelayPacket() *SuicideDelayPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"suicideDelayMS",
	}

	var id int32 = -911983090

	return &SuicideDelayPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-911983090, "SuicideDelay", func() packets.Packet {
		return NewSuicideDelayPacket()
	})
}
