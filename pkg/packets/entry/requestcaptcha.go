package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Request a captcha
type RequestCaptchaPacket struct {
	packets.BasePacket
}

func NewRequestCaptchaPacket() *RequestCaptchaPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"type_",
	}

	var id int32 = -349828108

	return &RequestCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-349828108, "RequestCaptchaPacket", func() packets.Packet {
		return NewRequestCaptchaPacket()
	})
}
