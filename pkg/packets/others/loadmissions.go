package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// User requests to load their missions
type LoadMissionsPacket struct {
	packets.BasePacket
}

func NewLoadMissionsPacket() *LoadMissionsPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoadMissionsID

	return &LoadMissionsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadMissionsID, "LoadMissions", func() packets.Packet {
		return NewLoadMissionsPacket()
	})
}
