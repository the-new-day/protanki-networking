package packets

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"maps"
	"slices"
)

// Boolshortern always returns an empty map.
// It is used to pass to Packet.UnwrapValues to say that the object is bool-shortened.
func Boolshortern() map[string]any {
	return make(map[string]any)
}

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
		id:         packet.ID(),
		codecs:     slices.Clone(packet.codecs),
		attributes: slices.Clone(packet.attributes),
		rawData:    slices.Clone(packet.rawData),
		objects:    slices.Clone(packet.objects),
		object:     maps.Clone(packet.object),
	}
}

// Compress comresses the data using raw DEFLATE algorithm.
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create compressor: %w", err)
	}
	defer w.Close()

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("compression failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize compression: %w", err)
	}

	return buf.Bytes(), nil
}

// Decompress decompresses data using raw DEFLATE algorithm.
func Decompress(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}
	return data, nil
}
