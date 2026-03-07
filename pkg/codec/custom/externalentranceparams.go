package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
)

type ExternalEntranceParamsCodec struct {
	CustomCodec
}

func NewExternalEntranceParamsCodec() *ExternalEntranceParamsCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "authUrl", complex.NewStringCodec())
	AddField(customCodec, "SnId", complex.NewStringCodec())

	return &ExternalEntranceParamsCodec{CustomCodec: *customCodec}
}
