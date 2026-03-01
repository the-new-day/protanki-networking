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

	var id int32 = packets.EndRespFantomID

	return &EndRespFantomPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.EndRespFantomID, "EndRespFantom", func() packets.Packet {
		return NewEndRespFantomPacket()
	})
}
