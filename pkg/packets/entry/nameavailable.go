package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Said name is available for registration
type NameAvailablePacket struct {
	packets.BasePacket
}

func NewNameAvailablePacket() *NameAvailablePacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -706679202

	return &NameAvailablePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-706679202, "NameAvailablePacket", func() packets.Packet {
		return NewNameAvailablePacket()
	})
}
