package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ChangeMissionID

	return &ChangeMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ChangeMissionID, "ChangeMission", func() packets.Packet {
		return NewChangeMissionPacket()
	})
}
