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

	var id int32 = -867767128

	return &AcceptMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-867767128, "AcceptMissionPacket", func() packets.Packet {
		return NewAcceptMissionPacket()
	})
}
