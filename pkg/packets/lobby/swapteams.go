package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Swaps teams in a battle
type SwapTeamsPacket struct {
	packets.BasePacket
}

func NewSwapTeamsPacket() *SwapTeamsPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = packets.SwapTeamsID

	return &SwapTeamsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SwapTeamsID, "SwapTeams", func() packets.Packet {
		return NewSwapTeamsPacket()
	})
}
