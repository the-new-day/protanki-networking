package custom

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/internal/codecs"
	"github.com/the-new-day/probogo/internal/codecs/primitive"
)

// Base codec for complex objects with multiple fields.
type CustomCodec struct {
	attributes   []string
	codecs       []codecs.Codec
	boolshortern bool
}

func NewCustomCodec(boolshortern bool) *CustomCodec {
	return &CustomCodec{
		attributes:   make([]string, 0),
		codecs:       make([]codecs.Codec, 0),
		boolshortern: boolshortern,
	}
}

func (c *CustomCodec) AddField(name string, codec codecs.Codec) {
	c.attributes = append(c.attributes, name)
	c.codecs = append(c.codecs, codec)
}

// Decodes a map from the buffer according to the defined fields.
func (c *CustomCodec) Decode(buf *bytes.Buffer) (map[string]any, error) {
	result := make(map[string]any)

	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty, err := boolCodec.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("CustomCodec: failed to decode empty flag: %w", err)
		}
		if isEmpty {
			return result, nil
		}
	}

	for i, attr := range c.attributes {
		val, err := c.codecs[i].Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("CustomCodec: failed to decode %s: %w", attr, err)
		}
		result[attr] = val
	}

	return result, nil
}

// Encodes a map to the buffer according to the defined fields.
func (c *CustomCodec) Encode(value map[string]any, buf *bytes.Buffer) (int, error) {
	totalBytes := 0

	if c.boolshortern {
		boolCodec := &primitive.BoolCodec{}
		isEmpty := len(value) == 0
		n, err := boolCodec.Encode(isEmpty, buf)
		totalBytes += n
		if err != nil {
			return totalBytes, fmt.Errorf("CustomCodec: failed to encode empty flag: %w", err)
		}

		if isEmpty {
			return totalBytes, nil
		}
	}

	for i, attr := range c.attributes {
		rawVal, ok := value[attr]
		if !ok {
			return totalBytes, fmt.Errorf("CustomCodec: missing attribute %q", attr)
		}

		n, err := c.codecs[i].Encode(rawVal, buf)
		totalBytes += n
		if err != nil {
			return totalBytes, fmt.Errorf("CustomCodec: failed to encode %s: %w", attr, err)
		}
	}

	return totalBytes, nil
}
