package shop

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client uses a promocode
type SendPromocodePacket struct {
	packets.BasePacket
}

func NewSendPromocodePacket() *SendPromocodePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"promocode",
	}

	var id int32 = -511004908

	return &SendPromocodePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
