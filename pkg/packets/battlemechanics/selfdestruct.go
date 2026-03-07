package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sends a self destruct request to the server
type SelfDestructPacket struct {
	packets.BasePacket
}

func NewSelfDestructPacket() *SelfDestructPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.SelfDestructID

	return &SelfDestructPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SelfDestructID, "SelfDestruct", func() packets.Packet {
		return NewSelfDestructPacket()
	})
}
