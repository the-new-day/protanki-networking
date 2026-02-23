package factory

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

// VectorCodec encodes/decodes a slice of elements of type T.
//
// Wire format:
//   - If boolshortern enabled: 1 byte empty flag (1 = empty, 0 = not empty)
//   - 4 bytes length (int32)
//   - N elements encoded with elementCodec (TypedCodec[T])
type VectorCodec[T any] struct {
	elementCodec codec.TypedCodec[T]
	boolshortern bool
}

func NewVectorCodec[T any](elemCodec codec.TypedCodec[T], boolshortern bool) *VectorCodec[T] {
	return &VectorCodec[T]{
		elementCodec: elemCodec,
		boolshortern: boolshortern,
	}
}

func (c *VectorCodec[T]) Decode(buf *bytes.Buffer) ([]T, error) {
	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty, err := boolCodec.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("VectorCodec: failed to decode empty flag: %w", err)
		}
		if isEmpty {
			return []T{}, nil
		}
	}

	intCodec := &primitive.IntCodec{}
	length, err := intCodec.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("VectorCodec: failed to decode length: %w", err)
	}

	result := make([]T, length)
	for i := range length {
		val, err := c.elementCodec.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("VectorCodec: failed to decode element at index %d: %w", i, err)
		}
		result[i] = val
	}

	return result, nil
}

func (c *VectorCodec[T]) Encode(value []T, buf *bytes.Buffer) (int, error) {
	totalBytes := 0

	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty := len(value) == 0
		n, err := boolCodec.Encode(isEmpty, buf)
		if err != nil {
			return totalBytes, fmt.Errorf("VectorCodec: failed to encode empty flag: %w", err)
		}
		totalBytes += n

		if isEmpty {
			return totalBytes, nil
		}
	}

	intCodec := &primitive.IntCodec{}
	n, err := intCodec.Encode(int32(len(value)), buf)
	if err != nil {
		return totalBytes, fmt.Errorf("VectorCodec: failed to encode length: %w", err)
	}
	totalBytes += n

	for i, item := range value {
		n, err := c.elementCodec.Encode(item, buf)
		if err != nil {
			return totalBytes, fmt.Errorf("VectorCodec: failed to encode element at index %d: %w", i, err)
		}
		totalBytes += n
	}

	return totalBytes, nil
}
