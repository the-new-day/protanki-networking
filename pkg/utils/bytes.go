package utils

import (
	"fmt"
	"io"
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
	if len(data) < 2*n {
		return fmt.Sprintf("% x", data)
	}
	return fmt.Sprintf("% x ... % x", data[:n], data[len(data)-n:])
}
