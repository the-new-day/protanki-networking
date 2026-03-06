package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads the player
type LoadFriendsListPacket struct {
	packets.BasePacket
}

func NewLoadFriendsListPacket() *LoadFriendsListPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
	}

	attributes := []string{
		"accepted",
		"newAccepted",
		"incoming",
		"newIncoming",
		"outgoing",
	}

	var id int32 = packets.LoadFriendsListID

	return &LoadFriendsListPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadFriendsListID, "LoadFriendsList", func() packets.Packet {
		return NewLoadFriendsListPacket()
	})
}
