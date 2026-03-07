package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type ChatMessageCodec struct {
	CustomCodec
}

func NewChatMessageCodec() *ChatMessageCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "authorStatus", NewUserStatusCodec())
	AddField(customCodec, "systemMessage", &primitive.BoolCodec{})
	AddField(customCodec, "targetStatus", NewUserStatusCodec())
	AddField(customCodec, "text", complex.NewStringCodec())
	AddField(customCodec, "warning", &primitive.BoolCodec{})

	return &ChatMessageCodec{CustomCodec: *customCodec}
}
