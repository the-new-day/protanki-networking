package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Open Friends List
type OpenFriendsListPacket struct {
	packets.BasePacket
}

func NewOpenFriendsListPacket() *OpenFriendsListPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 1441234714

	return &OpenFriendsListPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1441234714, "OpenFriendsListPacket", func() packets.Packet {
		return NewOpenFriendsListPacket()
	})
}
