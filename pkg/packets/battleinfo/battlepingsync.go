package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 2074243318

	return &BattlePingSyncPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(2074243318, "BattlePingSync", func() packets.Packet {
		return NewBattlePingSyncPacket()
	})
}
