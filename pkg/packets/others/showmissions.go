package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Show the list of missions currently available to the player
type ShowMissionsPacket struct {
	packets.BasePacket
}

func NewShowMissionsPacket() *ShowMissionsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewMissionCodec(), false)),
		codec.Wrap(custom.NewMissionStreakCodec()),
	}

	attributes := []string{
		"missions",
		"weeklyStreakInfo",
	}

	var id int32 = packets.ShowMissionsID

	return &ShowMissionsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ShowMissionsID, "ShowMissions", func() packets.Packet {
		return NewShowMissionsPacket()
	})
}
