package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Send Respawn(fantom) Packet
type SendRespawnPacket struct {
	packets.BasePacket
}

func NewSendRespawnPacket() *SendRespawnPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.SendRespawnID

	return &SendRespawnPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SendRespawnID, "SendRespawn", func() packets.Packet {
		return NewSendRespawnPacket()
	})
}
