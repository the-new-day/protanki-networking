package entry

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads the player
type LoadFriendsListPacket struct {
	packets.BasePacket
}

func NewLoadFriendsListPacket() *LoadFriendsListPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewVectorStringCodec()),
	}

	attributes := []string{
		"accepted",
		"newAccepted",
		"incoming",
		"newIncoming",
		"outgoing",
	}

	var id int32 = 1422563374

	return &LoadFriendsListPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
