package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Accept mission
type AcceptMissionPacket struct {
	packets.BasePacket
}

func NewAcceptMissionPacket() *AcceptMissionPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"missionId",
	}

	var id int32 = packets.AcceptMissionID

	return &AcceptMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.AcceptMissionID, "AcceptMission", func() packets.Packet {
		return NewAcceptMissionPacket()
	})
}
