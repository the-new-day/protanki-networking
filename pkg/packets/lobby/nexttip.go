package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Next tip (no attributes)
type NextTipPacket struct {
	packets.BasePacket
}

func NewNextTipPacket() *NextTipPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.NextTipID

	return &NextTipPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.NextTipID, "NextTip", func() packets.Packet {
		return NewNextTipPacket()
	})
}
