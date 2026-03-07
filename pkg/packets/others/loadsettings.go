package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Load settings
type LoadSettingsPacket struct {
	packets.BasePacket
}

func NewLoadSettingsPacket() *LoadSettingsPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoadSettingsID

	return &LoadSettingsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadSettingsID, "LoadSettings", func() packets.Packet {
		return NewLoadSettingsPacket()
	})
}
