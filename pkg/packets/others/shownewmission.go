package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Show the new mission that was previously changed
type ShowNewMissionPacket struct {
	packets.BasePacket
}

func NewShowNewMissionPacket() *ShowNewMissionPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(custom.NewMissionCodec()),
	}

	attributes := []string{
		"missionId",
		"mission",
	}

	var id int32 = packets.ShowNewMissionID

	return &ShowNewMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShowNewMissionID, "ShowMission", func() packets.Packet {
		return NewShowNewMissionPacket()
	})
}
