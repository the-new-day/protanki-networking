package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
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
