package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ChangeFreeMissionID

	return &ChangeFreeMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ChangeFreeMissionID, "ChangeFreeMission", func() packets.Packet {
		return NewChangeFreeMissionPacket()
	})
}
