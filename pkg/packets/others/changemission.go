package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
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
