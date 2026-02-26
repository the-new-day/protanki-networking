package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load battle info
type LoadLobbyPacket struct {
	packets.BasePacket
}

func NewLoadLobbyPacket() *LoadLobbyPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 1452181070

	return &LoadLobbyPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
