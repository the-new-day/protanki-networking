package complex

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

func TestStringCodec(t *testing.T) {
	codec := &StringCodec{}

	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"short", "Hello"},
		{"russian", "Привет"},
		{"with spaces", "Hello World"},
		{"special chars", "!@#$%^&*()"},
		{"unicode", "Go 世界"},
		{"very long", string(make([]byte, 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Encode
			n, err := codec.Encode(tt.value, buf)
			assert.NoError(t, err)

			// Verify size
			if tt.value == "" {
				assert.Equal(t, 1, n) // only empty flag
			} else {
				assert.Equal(t, 1+4+len(tt.value), n) // flag + length + data
			}

			// Decode
			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestStringCodec_InvalidData(t *testing.T) {
	codec := &StringCodec{}
	buf := &bytes.Buffer{}

	// Write only empty flag (0 = not empty) but no length
	buf.WriteByte(0)

	_, err := codec.Decode(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode length")
}

func TestDoubleIntCodec(t *testing.T) {
	codec := NewDoubleIntCodec("first", "second")

	tests := []struct {
		name   string
		first  int32
		second int32
	}{
		{"zero zero", 0, 0},
		{"positive positive", 42, 100},
		{"negative positive", -42, 55},
		{"max min", 2147483647, -2147483648},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			data := map[string]int32{
				"first":  tt.first,
				"second": tt.second,
			}

			n, err := codec.Encode(data, buf)
			assert.NoError(t, err)
			assert.Equal(t, 1+4+4, n) // boolshortern(1) + first(4) + second(4)

			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.first, result["first"])
			assert.Equal(t, tt.second, result["second"])
		})
	}
}

func TestDoubleIntCodec_Empty(t *testing.T) {
	codec := NewDoubleIntCodec("first", "second")
	buf := &bytes.Buffer{}

	// Encode empty (should still write boolshortern flag = false because not empty)
	data := map[string]int32{
		"first":  int32(0),
		"second": int32(0),
	}
	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1+4+4, n) // boolshortern(1) + first(4) + second(4)

	// Empty map should have boolshortern = true
	buf.Reset()
	n, err = codec.Encode(map[string]int32{}, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n) // only boolshortern flag
}

func TestVector3DCodec(t *testing.T) {
	codec := NewVector3DCodec()

	tests := []struct {
		name    string
		x, y, z float32
	}{
		{"zero", 0, 0, 0},
		{"positive", 1.5, 2.3, 3.7},
		{"negative", -1.0, -2.0, -3.0},
		{"mixed", 1.0, -2.5, 3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			data := map[string]float32{
				"x": tt.x,
				"y": tt.y,
				"z": tt.z,
			}

			n, err := codec.Encode(data, buf)
			assert.NoError(t, err)
			assert.Equal(t, 1+4+4+4, n) // boolshortern(1) + x,y,z(4 each)

			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.InDelta(t, tt.x, result["x"], 0.0001)
			assert.InDelta(t, tt.y, result["y"], 0.0001)
			assert.InDelta(t, tt.z, result["z"], 0.0001)
		})
	}
}

func TestVectorVector3DCodec(t *testing.T) {
	codec := NewVectorVector3DCodec()

	tests := []struct {
		name    string
		vectors []map[string]float32
	}{
		{
			name:    "empty",
			vectors: []map[string]float32{},
		},
		{
			name: "single",
			vectors: []map[string]float32{
				{"x": float32(1.0), "y": float32(2.0), "z": float32(3.0)},
			},
		},
		{
			name: "multiple",
			vectors: []map[string]float32{
				{"x": float32(1.0), "y": float32(2.0), "z": float32(3.0)},
				{"x": float32(4.0), "y": float32(5.0), "z": float32(6.0)},
				{"x": float32(7.0), "y": float32(8.0), "z": float32(9.0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Encode
			n, err := codec.Encode(tt.vectors, buf)
			assert.NoError(t, err)

			// Verify size
			if len(tt.vectors) == 0 {
				assert.Equal(t, 1, n) // only empty flag
			} else {
				expectedSize := 1 + 4 // vector empty flag + length
				for range tt.vectors {
					expectedSize += 1 + 4 + 4 + 4 // vec3 flag + x,y,z
				}
				assert.Equal(t, expectedSize, n)
			}

			// Decode
			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.vectors), len(result))

			for i, vec := range result {
				expected := tt.vectors[i]
				assert.InDelta(t, expected["x"], vec["x"], 0.0001)
				assert.InDelta(t, expected["y"], vec["y"], 0.0001)
				assert.InDelta(t, expected["z"], vec["z"], 0.0001)
			}
		})
	}
}

func TestVectorStringCodec(t *testing.T) {
	codec := NewVectorStringCodec()

	tests := []struct {
		name    string
		strings []string
	}{
		{
			name:    "empty",
			strings: []string{},
		},
		{
			name:    "single",
			strings: []string{"hello"},
		},
		{
			name:    "multiple",
			strings: []string{"hello", "world", "!"},
		},
		{
			name:    "with empty",
			strings: []string{"", "non-empty", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			n, err := codec.Encode(tt.strings, buf)
			assert.NoError(t, err)

			if len(tt.strings) == 0 {
				assert.Equal(t, 1, n) // only empty flag
			}

			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.strings, result)
		})
	}
}

func TestVectorShortCodec(t *testing.T) {
	codec := NewVectorShortCodec()

	tests := []struct {
		name   string
		values []int16
	}{
		{
			name:   "empty",
			values: []int16{},
		},
		{
			name:   "single",
			values: []int16{42},
		},
		{
			name:   "multiple",
			values: []int16{1, 2, 3, 4, 5},
		},
		{
			name:   "negative",
			values: []int16{-1, -2, -3},
		},
		{
			name:   "mixed",
			values: []int16{-32768, 0, 32767},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			n, err := codec.Encode(tt.values, buf)
			assert.NoError(t, err)

			if len(tt.values) == 0 {
				assert.Equal(t, 1, n) // only empty flag
			} else {
				assert.Equal(t, 1+4+len(tt.values)*2, n) // flag + length + elements
			}

			result, err := codec.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.values, result)
		})
	}
}

func TestErrorCases(t *testing.T) {
	t.Run("StringCodec - invalid empty flag", func(t *testing.T) {
		codec := &StringCodec{}
		buf := &bytes.Buffer{}

		// No data at all
		_, err := codec.Decode(buf)
		assert.Error(t, err)
	})

	t.Run("VectorCodec - invalid data", func(t *testing.T) {
		vecCodec := NewVectorStringCodec()
		buf := &bytes.Buffer{}

		// Write only empty flag = false, but no length
		boolCodec := &primitive.BoolCodec{}
		boolCodec.Encode(false, buf)

		_, err := vecCodec.Decode(buf)
		assert.Error(t, err)
	})

	t.Run("MultiCodec - missing attribute", func(t *testing.T) {
		codec := NewDoubleIntCodec("first", "second")
		buf := &bytes.Buffer{}

		// Missing "second"
		data := map[string]int32{
			"first": 42,
		}

		_, err := codec.Encode(data, buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing attribute")
	})
}

func TestVectorVector3DCodec_WithBoolshortern(t *testing.T) {
	codec := NewVectorVector3DCodec()
	buf := &bytes.Buffer{}

	// Test with empty slice
	emptySlice := []map[string]float32{}
	n, err := codec.Encode(emptySlice, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n) // only vector empty flag

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Test with slice containing empty vectors
	buf.Reset()
	vectorsWithEmpty := []map[string]float32{
		{}, // empty vector (boolshortern = true)
		{"x": float32(1.0), "y": float32(2.0), "z": float32(3.0)},
	}

	n, err = codec.Encode(vectorsWithEmpty, buf)
	assert.NoError(t, err)

	result, err = codec.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Empty(t, result[0])
	assert.NotEmpty(t, result[1])
}

func TestJsonCodec(t *testing.T) {
	codec := &JsonCodec{}

	tests := []struct {
		name  string
		value map[string]any
	}{
		{
			name:  "empty",
			value: map[string]any{},
		},
		{
			name: "simple",
			value: map[string]any{
				"name": "Alice",
				"age":  30,
			},
		},
		{
			name: "nested",
			value: map[string]any{
				"user": map[string]any{
					"id":   123,
					"name": "Bob",
				},
				"tags": []string{"player", "premium"},
			},
		},
		{
			name: "mixed types",
			value: map[string]any{
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"string": "hello",
				"null":   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Encode
			n, err := codec.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.True(t, n > 0)

			// Decode
			result, err := codec.Decode(buf)
			assert.NoError(t, err)

			expected := tt.value

			assert.True(t, compareJSONValues(t, expected, result), "JSON values don't match")
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func compareJSONValues(t *testing.T, expected, actual any) bool {
	switch exp := expected.(type) {
	case map[string]any:
		act, ok := actual.(map[string]any)
		if !ok {
			t.Errorf("expected map, got %T", actual)
			return false
		}
		if len(exp) != len(act) {
			t.Errorf("map length mismatch: %d vs %d", len(exp), len(act))
			return false
		}
		for k, v := range exp {
			if !compareJSONValues(t, v, act[k]) {
				t.Errorf("key %q mismatch", k)
				return false
			}
		}
		return true

	case []any:
		act, ok := actual.([]any)
		if !ok {
			t.Errorf("expected slice, got %T", actual)
			return false
		}
		if len(exp) != len(act) {
			return false
		}
		for i, v := range exp {
			if !compareJSONValues(t, v, act[i]) {
				return false
			}
		}
		return true

	case []string:
		act, ok := actual.([]any)
		if !ok {
			t.Errorf("expected []any for string slice, got %T", actual)
			return false
		}
		if len(exp) != len(act) {
			return false
		}
		for i, v := range exp {
			if v != act[i].(string) {
				return false
			}
		}
		return true

	case int:
		act, ok := actual.(float64)
		if !ok {
			t.Errorf("expected float64 for int, got %T", actual)
			return false
		}
		return float64(exp) == act

	case float64:
		act, ok := actual.(float64)
		if !ok {
			return false
		}
		return exp == act

	case string:
		act, ok := actual.(string)
		if !ok {
			return false
		}
		return exp == act

	case bool:
		act, ok := actual.(bool)
		if !ok {
			return false
		}
		return exp == act

	case nil:
		return actual == nil

	default:
		t.Errorf("unexpected type: %T", exp)
		return false
	}
}

func TestJsonCodec_EmptyString(t *testing.T) {
	codec := &JsonCodec{}
	buf := &bytes.Buffer{}

	// Encode empty map
	n, err := codec.Encode(map[string]any{}, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1+4+2, n) // flag(1) + length(4) + "{}"(2)

	// Decode
	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestJsonCodec_InvalidJson(t *testing.T) {
	codec := &JsonCodec{}
	buf := &bytes.Buffer{}

	// Write invalid JSON string
	stringCodec := &StringCodec{}
	stringCodec.Encode("{invalid json", buf)

	_, err := codec.Decode(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}
