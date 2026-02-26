package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Buy from shop
type ChangeLocationPacket struct {
	packets.BasePacket
}

func NewChangeLocationPacket() *ChangeLocationPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"location_abbreviation",
	}

	var id int32 = 921004371

	return &ChangeLocationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
