package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Open Friends List
type OpenFriendsListPacket struct {
	packets.BasePacket
}

func NewOpenFriendsListPacket() *OpenFriendsListPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.OpenFriendsListID

	return &OpenFriendsListPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.OpenFriendsListID, "OpenFriendsList", func() packets.Packet {
		return NewOpenFriendsListPacket()
	})
}
