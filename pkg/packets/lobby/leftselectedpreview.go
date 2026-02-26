package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A player has left the selected battle, from the preview screen
type LeftSelectedPreviewPacket struct {
	packets.BasePacket
}

func NewLeftSelectedPreviewPacket() *LeftSelectedPreviewPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
		"username",
	}

	var id int32 = 1924874982

	return &LeftSelectedPreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
