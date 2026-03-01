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

	var id int32 = packets.LoadReferralID

	return &LoadReferralPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadReferralID, "LoadReferral", func() packets.Packet {
		return NewLoadReferralPacket()
	})
}
