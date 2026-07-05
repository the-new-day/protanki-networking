package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Load battle info (C2S)
type LoadLobbyPacket struct {
	packets.BasePacket
}

func NewLoadLobbyPacket() *LoadLobbyPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoadLobbyID

	return &LoadLobbyPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadLobbyID, "LoadLobby", func() packets.Packet {
		return NewLoadLobbyPacket()
	})
}
