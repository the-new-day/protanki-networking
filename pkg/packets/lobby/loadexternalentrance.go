package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads external entrance data, containing a vector of social network entrance params.
type LoadExternalEntrancePacket struct {
	packets.BasePacket
}

func NewLoadExternalEntrancePacket() *LoadExternalEntrancePacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewExternalEntranceParamsCodec(), false)),
	}

	attributes := []string{
		"params",
	}

	var id int32 = packets.LoadExternalEntranceID

	return &LoadExternalEntrancePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadExternalEntranceID, "LoadExternalEntrance", func() packets.Packet {
		return NewLoadExternalEntrancePacket()
	})
}
