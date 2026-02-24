package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Send respawn (fantom) packet.
type SendRespawnPacket struct {
	packets.BasePacket
}

func NewSendRespawnPacket() *SendRespawnPacket {
	codecs := []codec.Codec{}
	attributes := []string{}

	var id int32 = -1378839846

	return &SendRespawnPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
