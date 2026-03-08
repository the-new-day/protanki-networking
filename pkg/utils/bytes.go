package utils

import (
	"fmt"
	"io"
	"strings"
)

// ReadBytes reads exactly n bytes from r. It returns a slice of length n
// containing the bytes, or an error if fewer than n bytes are read.
func ReadBytes(n int, r io.Reader) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

// ShortView returns a compact hexadecimal representation of data.
// If len(data) < 2*n, it returns the whole slice in hex.
// Otherwise, it returns the first n and last n bytes, separated by " ... ".
//
// Example:
//
//	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
//	view := ShortView(data, 2) // 00 01 ... 04 05
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
