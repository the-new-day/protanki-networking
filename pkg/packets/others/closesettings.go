package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -731115522

	return &CloseSettingsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
