package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load settings
type LoadSettingsPacket struct {
	packets.BasePacket
}

func NewLoadSettingsPacket() *LoadSettingsPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 850220815

	return &LoadSettingsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(850220815, "LoadSettingsPacket", func() packets.Packet {
        return NewLoadSettingsPacket()
    })
}
