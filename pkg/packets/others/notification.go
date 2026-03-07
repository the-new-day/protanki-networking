package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.SettingsNotificationID

	return &SettingsNotificationPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SettingsNotificationID, "SettingsNotification", func() packets.Packet {
		return NewSettingsNotificationPacket()
	})
}
