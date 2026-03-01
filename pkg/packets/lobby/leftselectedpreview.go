package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
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

	var id int32 = packets.LeftSelectedPreviewID

	return &LeftSelectedPreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LeftSelectedPreviewID, "LeftSelectedPreview", func() packets.Packet {
		return NewLeftSelectedPreviewPacket()
	})
}
