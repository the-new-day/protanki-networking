package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Attempts to sync ping information with the server
type BattlePingSyncPacket struct {
	packets.BasePacket
}

func NewBattlePingSyncPacket() *BattlePingSyncPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewDoubleIntCodec("clientTime", "serverSessionTime")),
	}

	attributes := []string{
		"latencyInfo",
	}

	var id int32 = packets.BattlePingSyncID

	return &BattlePingSyncPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattlePingSyncID, "BattlePingSync", func() packets.Packet {
		return NewBattlePingSyncPacket()
	})
}
