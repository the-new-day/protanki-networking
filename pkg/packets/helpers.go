package packets

import "fmt"

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
