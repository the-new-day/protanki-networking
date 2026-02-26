package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sends server details about a railgun shot that has just started to release
type RailgunShotInitOutPacket struct {
	packets.BasePacket
}

func NewRailgunShotInitOutPacket() *RailgunShotInitOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -1759063234

	return &RailgunShotInitOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
