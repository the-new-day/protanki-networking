package shop

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
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

	var id int32 = packets.SendPromocodeID

	return &SendPromocodePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SendPromocodeID, "SendPromocode", func() packets.Packet {
		return NewSendPromocodePacket()
	})
}
