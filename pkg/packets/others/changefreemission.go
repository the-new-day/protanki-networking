package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Change free mission
type ChangeFreeMissionPacket struct {
	packets.BasePacket
}

func NewChangeFreeMissionPacket() *ChangeFreeMissionPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"missionId",
	}

	var id int32 = 326032325

	return &ChangeFreeMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(326032325, "ChangeFreeMission", func() packets.Packet {
		return NewChangeFreeMissionPacket()
	})
}
