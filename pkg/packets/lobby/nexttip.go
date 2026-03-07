package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Client asks for the next tip on the loading screen (no attributes)
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
