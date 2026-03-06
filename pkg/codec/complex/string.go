package complex

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/utils"
)

type StringCodec struct {
	primitive.BoolCodec
	primitive.IntCodec
}

func NewStringCodec() *StringCodec {
	return &StringCodec{}
}

// 1 Byte - Is string empty? (If empty, everything else is ignored)
// 4 Bytes - String length (bytes)
// Remaining Bytes - String value

func (c *StringCodec) Decode(buf *bytes.Buffer) (string, error) {
	isEmpty, err := c.BoolCodec.Decode(buf)
	if err != nil {
		return "", fmt.Errorf("StringCodec: failed to decode empty flag: %w", err)
	}
	if isEmpty {
		return "", nil
	}

	length, err := c.IntCodec.Decode(buf)
	if err != nil {
		return "", fmt.Errorf("StringCodec: failed to decode length: %w", err)
	}

	if length < 0 {
		return "", fmt.Errorf("StringCodec: invalid negative length: %d", length)
	}

	stringBytes, err := utils.ReadBytes(int(length), buf)
	if err != nil {
		return "", fmt.Errorf("StringCodec: failed to read %d bytes of content: %w", length, err)
	}

	return string(stringBytes), nil
}

func (c *StringCodec) Encode(value string, buf *bytes.Buffer) (int, error) {
	rawBytes := []byte(value)
	isEmpty := len(value) == 0
	total := 0

	n, err := c.BoolCodec.Encode(isEmpty, buf)
	if err != nil {
		return n, fmt.Errorf("StringCodec: failed to encode empty flag: %w", err)
	}
	total += n

	if isEmpty {
		return total, nil
	}

	n, err = c.IntCodec.Encode(int32(len(rawBytes)), buf)
	if err != nil {
		return total, fmt.Errorf("StringCodec: failed to encode length: %w", err)
	}
	total += n

	n, err = buf.Write(rawBytes)
	if err != nil {
		return total, fmt.Errorf("StringCodec: failed to write content: %w", err)
	}
	total += n

	return total, nil
}
