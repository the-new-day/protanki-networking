package utils

import (
	"fmt"
	"io"
)

func ReadBytes(n int, r io.Reader) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func ShortView(data []byte, n int, m int) string {
	if len(data) < n {
		return fmt.Sprintf("{% x}", data)
	}
	return fmt.Sprintf("{% x ... % x}", data[:m], data[len(data)-m:])
}
