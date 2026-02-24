package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type ChatMessageCodec struct {
	CustomCodec
}

func NewChatMessageCodec() *ChatMessageCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "authorStatus", &UserStatusCodec{})
	AddField(customCodec, "systemMessage", &primitive.BoolCodec{})
	AddField(customCodec, "targetStatus", &UserStatusCodec{})
	AddField(customCodec, "message", &complex.StringCodec{})
	AddField(customCodec, "warning", &primitive.BoolCodec{})

	return &ChatMessageCodec{CustomCodec: *customCodec}
}
