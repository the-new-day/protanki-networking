package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Initiates a shaft scope shot
type ShaftScopeInitOutPacket struct {
	packets.BasePacket
}

func NewShaftScopeInitOutPacket() *ShaftScopeInitOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -367760678

	return &ShaftScopeInitOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
