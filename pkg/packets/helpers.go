package packets

import (
	"fmt"
	"maps"
	"slices"
)

// Attr returns value of the attribute with the given name,
// converted to the type T. Calls Packet.Attr for retrieving the value.
// Returns converted value or panics on the type mismatch or if the attribute doesn't exist.
func Attr[T any](name string, packet Packet) T {
	rawValue := packet.Attr(name)
	var zero T

	value, ok := rawValue.(T)
	if !ok {
		panic(fmt.Sprintf(
			"Attr: packet %d attribute %q type mismatch: want %T, got %T",
			packet.ID(), name, zero, rawValue,
		))
	}
	return value
}

// Clone returns copy of a packet.
// It uses slices.Clone and maps.Clone for each slices and map fields,
// so the behaviour for them is the same.
func Clone(packet *BasePacket) Packet {
	return &BasePacket{
		id:               packet.ID(),
		codecs:           slices.Clone(packet.codecs),
		attributes:       slices.Clone(packet.attributes),
		encryptedRawData: slices.Clone(packet.encryptedRawData),
		rawData:          slices.Clone(packet.rawData),
		objects:          slices.Clone(packet.objects),
		object:           maps.Clone(packet.object),
	}
}
