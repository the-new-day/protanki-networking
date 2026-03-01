package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A mission was completed and is available to claim
type CompletedMissionPacket struct {
	packets.BasePacket
}

func NewCompletedMissionPacket() *CompletedMissionPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 1579425801

	return &CompletedMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1579425801, "CompletedMission", func() packets.Packet {
		return NewCompletedMissionPacket()
	})
}
