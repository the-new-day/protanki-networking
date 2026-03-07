package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// The captcha was incorrect, a new one is sent
type WrongNewCaptchaPacket struct {
	packets.BasePacket
}

func NewWrongNewCaptchaPacket() *WrongNewCaptchaPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(multiple.NewVectorCodec(&primitive.ByteCodec{}, false)),
	}

	attributes := []string{
		"type_",
		"imagedata",
	}

	var id int32 = packets.WrongNewCaptchaID

	return &WrongNewCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.WrongNewCaptchaID, "WrongCaptcha", func() packets.Packet {
		return NewWrongNewCaptchaPacket()
	})
}
