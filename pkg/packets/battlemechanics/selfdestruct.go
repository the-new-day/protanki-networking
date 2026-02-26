package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends a self destruct request to the server
type SelfDestructPacket struct {
	packets.BasePacket
}

func NewSelfDestructPacket() *SelfDestructPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 988664577

	return &SelfDestructPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
