package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ShaftScopeInitOutID

	return &ShaftScopeInitOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShaftScopeInitOutID, "ShaftScopeInitOut", func() packets.Packet {
		return NewShaftScopeInitOutPacket()
	})
}
