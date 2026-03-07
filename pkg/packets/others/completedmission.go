package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// A mission was completed and is available to claim
type CompletedMissionPacket struct {
	packets.BasePacket
}

func NewCompletedMissionPacket() *CompletedMissionPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.CompletedMissionID

	return &CompletedMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CompletedMissionID, "CompletedMission", func() packets.Packet {
		return NewCompletedMissionPacket()
	})
}
