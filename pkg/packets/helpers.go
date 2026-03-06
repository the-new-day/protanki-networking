package packets

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
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
		id:         packet.ID(),
		codecs:     slices.Clone(packet.codecs),
		attributes: slices.Clone(packet.attributes),
		rawData:    slices.Clone(packet.rawData),
		objects:    slices.Clone(packet.objects),
		object:     maps.Clone(packet.object),
	}
}

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

// ShortView returns a hex representation of a byte slice, showing only
// the first and last n bytes if the slice is longer than 2n bytes.
// The result format is like: [54 6d 09 ... 1d f1 2e]
func ShortView(data []byte, n int) string {
	if len(data) <= 2*n {
		// Show all bytes
		var parts []string
		for _, b := range data {
			parts = append(parts, fmt.Sprintf("%02x", b))
		}
		return strings.Join(parts, " ")
	}

	// Show first n and last n bytes
	first := data[:n]
	last := data[len(data)-n:]

	var firstParts, lastParts []string
	for _, b := range first {
		firstParts = append(firstParts, fmt.Sprintf("%02x", b))
	}
	for _, b := range last {
		lastParts = append(lastParts, fmt.Sprintf("%02x", b))
	}

	return strings.Join(firstParts, " ") + " ... " + strings.Join(lastParts, " ")
}
