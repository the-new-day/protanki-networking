package packets

import "fmt"

func Attr[T any](name string, packet Packet) (T, error) {
	rawValue, err := packet.Attr(name)
	var zero T

	if err != nil {
		return zero, fmt.Errorf("Attr: %w", err)
	}

	value, ok := rawValue.(T)
	if !ok {
		return zero, fmt.Errorf(
			"Attr: packet %d attribute %q type mismatch: want %T, got %T",
			packet.ID(), name, zero, rawValue,
		)
	}
	return value, nil
}
