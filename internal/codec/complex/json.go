package complex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonCodec struct {
	StringCodec
}

func (c *JsonCodec) Decode(buf *bytes.Buffer) (map[string]any, error) {
	jsonStr, err := c.StringCodec.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("JsonCodec: failed to read string: %w", err)
	}

	if jsonStr == "" {
		return map[string]any{}, nil
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("JsonCodec: failed to parse: %w", err)
	}

	return result, nil
}

func (c *JsonCodec) Encode(value map[string]any, buf *bytes.Buffer) (int, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return 0, fmt.Errorf("JsonCodec: failed to marshal: %w", err)
	}

	return c.StringCodec.Encode(string(jsonBytes), buf)
}
