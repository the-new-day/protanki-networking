package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load garage
type LoadReferralPacket struct {
	packets.BasePacket
}

func NewLoadReferralPacket() *LoadReferralPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -169921234

	return &LoadReferralPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-169921234, "LoadReferralPacket", func() packets.Packet {
		return NewLoadReferralPacket()
	})
}
