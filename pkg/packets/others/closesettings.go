package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Close settings modal
type CloseSettingsPacket struct {
	packets.BasePacket
}

func NewCloseSettingsPacket() *CloseSettingsPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"close_state",
	}

	var id int32 = packets.CloseSettingsID

	return &CloseSettingsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CloseSettingsID, "CloseSettings", func() packets.Packet {
		return NewCloseSettingsPacket()
	})
}
