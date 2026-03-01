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

	var id int32 = packets.RequestCaptchaID

	return &RequestCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RequestCaptchaID, "RequestCaptcha", func() packets.Packet {
		return NewRequestCaptchaPacket()
	})
}
