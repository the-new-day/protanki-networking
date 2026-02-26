package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// User requests to load their missions
type LoadMissionsPacket struct {
	packets.BasePacket
}

func NewLoadMissionsPacket() *LoadMissionsPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 1227293080

	return &LoadMissionsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
