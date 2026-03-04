package complex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonCodec struct {
	StringCodec
}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Decode(buf *bytes.Buffer) (map[string]any, error) {
	jsonStr, err := c.StringCodec.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("JsonCodec: failed to read string: %w", err)
	}

	if jsonStr == "" {
		return map[string]any{}, nil
	}

	var data any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("JsonCodec: failed to parse: %w", err)
	}

	// Handle both objects and arrays
	switch v := data.(type) {
	case map[string]any:
		return v, nil
	case []any:
		// Convert array to map with stringified indices
		result := make(map[string]any)
		for i, item := range v {
			result[fmt.Sprintf("%d", i)] = item
		}
		return result, nil
	default:
		return nil, fmt.Errorf("JsonCodec: expected object or array, got %T", v)
	}
}

func (c *JsonCodec) Encode(value map[string]any, buf *bytes.Buffer) (int, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return 0, fmt.Errorf("JsonCodec: failed to marshal: %w", err)
	}

	return c.StringCodec.Encode(string(jsonBytes), buf)
}
