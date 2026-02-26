package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sets client language
type SetClientLangPacket struct {
	packets.BasePacket
}

func NewSetClientLangPacket() *SetClientLangPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"lang",
	}

	var id int32 = -1864333717

	return &SetClientLangPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
