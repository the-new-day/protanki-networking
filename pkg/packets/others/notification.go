package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Check if the user has notification on
type SettingsNotificationPacket struct {
	packets.BasePacket
}

func NewSettingsNotificationPacket() *SettingsNotificationPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"notificationEnabled",
	}

	var id int32 = 1447082276

	return &SettingsNotificationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1447082276, "SettingsNotificationPacket", func() packets.Packet {
        return NewSettingsNotificationPacket()
    })
}
