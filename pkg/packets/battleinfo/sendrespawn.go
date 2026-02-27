package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Send Respawn(fantom) Packet
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


func init() {
    packets.Register(-1378839846, "SendRespawnPacket", func() packets.Packet {
        return NewSendRespawnPacket()
    })
}
