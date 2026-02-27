package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -994817471

	return &SwapTeamsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-994817471, "SwapTeamsPacket", func() packets.Packet {
        return NewSwapTeamsPacket()
    })
}
