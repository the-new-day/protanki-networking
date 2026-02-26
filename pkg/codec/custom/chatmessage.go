package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type ChatMessageCodec struct {
	CustomCodec
}

func NewChatMessageCodec() *ChatMessageCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "authorStatus", NewUserStatusCodec())
	AddField(customCodec, "systemMessage", &primitive.BoolCodec{})
	AddField(customCodec, "targetStatus", NewUserStatusCodec())
	AddField(customCodec, "message", complex.NewStringCodec())
	AddField(customCodec, "warning", &primitive.BoolCodec{})

	return &ChatMessageCodec{CustomCodec: *customCodec}
}
