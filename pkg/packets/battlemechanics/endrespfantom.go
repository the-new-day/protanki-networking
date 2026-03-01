package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client end the Fantom period to fully spawn
type EndRespFantomPacket struct {
	packets.BasePacket
}

func NewEndRespFantomPacket() *EndRespFantomPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 1178028365

	return &EndRespFantomPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1178028365, "EndRespFantom", func() packets.Packet {
		return NewEndRespFantomPacket()
	})
}
