package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Information about the ping of the player.
type BattlePingInfoPacket struct {
	packets.BasePacket
}

func NewBattlePingInfoPacket() *BattlePingInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewDoubleIntCodec("serverSessionTime", "clientPing")),
	}

	attributes := []string{
		"latencyInfo",
	}

	var id int32 = 34068208

	return &BattlePingInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(34068208, "BattlePingInfoPacket", func() packets.Packet {
		return NewBattlePingInfoPacket()
	})
}
