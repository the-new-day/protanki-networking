package turrets

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Firebird stops shooting
type FireEndOutPacket struct {
	packets.BasePacket
}

func NewFireEndOutPacket() *FireEndOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -1300958299

	return &FireEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
