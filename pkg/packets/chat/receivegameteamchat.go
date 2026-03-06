package chat

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Receives a message from the game team chat
type ReceiveGameTeamChat struct {
	packets.BasePacket
}

func NewBattleChatAddTeamMessagePacket() *ReceiveGameTeamChat {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"message",
		"team",
	}

	var id int32 = packets.BattleChatAddTeamMessageID

	return &ReceiveGameTeamChat{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleChatAddTeamMessageID, "BattleChatAddTeamMessage", func() packets.Packet {
		return NewBattleChatAddTeamMessagePacket()
	})
}
