package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -1266665816

	return &ShowNewMissionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
