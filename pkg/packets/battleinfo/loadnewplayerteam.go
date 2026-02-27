package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A new player has joined the team battle
type LoadNewPlayerTeamPacket struct {
	packets.BasePacket
}

func NewLoadNewPlayerTeamPacket() *LoadNewPlayerTeamPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserCodec(), false)),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"userinfos",
		"team",
	}

	var id int32 = 2040021062

	return &LoadNewPlayerTeamPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(2040021062, "LoadPlayerTeamPacket", func() packets.Packet {
        return NewLoadNewPlayerTeamPacket()
    })
}
