package multiple

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

// MultiCodec encodes/decodes multiple fields of the same type T.
//
// Wire format:
//   - If boolshortern enabled: 1 byte empty flag (1 = empty, 0 = not empty)
//   - N values encoded with elementCodec (TypedCodec[T]) in attributes order
type MultiCodec[T any] struct {
	attributes   []string
	elementCodec codec.TypedCodec[T]
	boolshortern bool
}

func NewMultiCodec[T any](attrs []string, elemCodec codec.TypedCodec[T], boolshortern bool) *MultiCodec[T] {
	return &MultiCodec[T]{
		attributes:   attrs,
		elementCodec: elemCodec,
		boolshortern: boolshortern,
	}
}

func (c *MultiCodec[T]) Decode(buf *bytes.Buffer) (map[string]T, error) {
	result := make(map[string]T)

	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty, err := boolCodec.Decode(buf)
		if err != nil {
			return nil, err
		}
		if isEmpty {
			return result, nil
		}
	}

	for _, attr := range c.attributes {
		val, err := c.elementCodec.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("MultiCodec: failed to decode %s: %w", attr, err)
		}
		result[attr] = val
	}

	return result, nil
}

func (c *MultiCodec[T]) Encode(value map[string]T, buf *bytes.Buffer) (int, error) {
	totalBytes := 0

	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty := len(value) == 0
		n, err := boolCodec.Encode(isEmpty, buf)
		if err != nil {
			return totalBytes, err
		}
		totalBytes += n

		if isEmpty {
			return totalBytes, nil
		}
	}

	for _, attr := range c.attributes {
		val, ok := value[attr]
		if !ok {
			return totalBytes, fmt.Errorf("MultiCodec: missing attribute %q", attr)
		}

		n, err := c.elementCodec.Encode(val, buf)
		if err != nil {
			return totalBytes, fmt.Errorf("MultiCodec: failed to encode %s: %w", attr, err)
		}
		totalBytes += n
	}

	return totalBytes, nil
}
