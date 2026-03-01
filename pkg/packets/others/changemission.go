package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Change mission
type ChangeMissionPacket struct {
	packets.BasePacket
}

func NewChangeMissionPacket() *ChangeMissionPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"missionId",
	}

	var id int32 = 1642608662

	return &ChangeMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1642608662, "ChangeMission", func() packets.Packet {
		return NewChangeMissionPacket()
	})
}
