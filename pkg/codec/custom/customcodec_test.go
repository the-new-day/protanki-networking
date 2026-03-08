package custom

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

func TestCustomCodec_EncodeDecode(t *testing.T) {
	tests := []struct {
		name            string
		setup           func() *CustomCodec
		input           map[string]any
		validate        func(t *testing.T, result map[string]any)
		expectedSize    int
		expectEncodeErr bool
	}{
		{
			name: "without boolshortern",
			setup: func() *CustomCodec {
				cc := NewCustomCodec(false)
				cc.AddField("name", codec.Wrap(complex.NewStringCodec()))
				cc.AddField("age", codec.Wrap(&primitive.IntCodec{}))
				cc.AddField("active", codec.Wrap(&primitive.BoolCodec{}))
				return cc
			},
			input: map[string]any{
				"name":   "Alice",
				"age":    int32(30),
				"active": true,
			},
			expectedSize: 1 + 4 + len("Alice") + 4 + 1,
			validate: func(t *testing.T, result map[string]any) {
				assert.Equal(t, "Alice", result["name"])
				assert.Equal(t, int32(30), result["age"])
				assert.Equal(t, true, result["active"])
			},
		},
		{
			name: "with boolshortern non empty",
			setup: func() *CustomCodec {
				cc := NewCustomCodec(true)
				cc.AddField("x", codec.Wrap(&primitive.FloatCodec{}))
				cc.AddField("y", codec.Wrap(&primitive.FloatCodec{}))
				return cc
			},
			input: map[string]any{
				"x": float32(1.5),
				"y": float32(2.7),
			},
			expectedSize: 1 + 4 + 4,
			validate: func(t *testing.T, result map[string]any) {
				assert.InDelta(t, float32(1.5), result["x"].(float32), 0.0001)
				assert.InDelta(t, float32(2.7), result["y"].(float32), 0.0001)
			},
		},
		{
			name: "boolshortern empty map",
			setup: func() *CustomCodec {
				cc := NewCustomCodec(true)
				cc.AddField("x", codec.Wrap(&primitive.FloatCodec{}))
				cc.AddField("y", codec.Wrap(&primitive.FloatCodec{}))
				return cc
			},
			input:        packets.Boolshortern(),
			expectedSize: 1,
			validate: func(t *testing.T, result map[string]any) {
				assert.Empty(t, result)
			},
		},
		{
			name: "boolshortern encode empty even with fields",
			setup: func() *CustomCodec {
				cc := NewCustomCodec(true)
				cc.AddField("a", codec.Wrap(&primitive.IntCodec{}))
				cc.AddField("b", codec.Wrap(&primitive.IntCodec{}))
				return cc
			},
			input:        packets.Boolshortern(),
			expectedSize: 1,
			validate: func(t *testing.T, result map[string]any) {
				assert.Empty(t, result)
			},
		},
		{
			name: "no fields with boolshortern",
			setup: func() *CustomCodec {
				return NewCustomCodec(true)
			},
			input:        map[string]any{"unused": "value"},
			expectedSize: 1,
			validate: func(t *testing.T, result map[string]any) {
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.setup()
			buf := &bytes.Buffer{}

			n, err := cc.Encode(tt.input, buf)

			if tt.expectEncodeErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.expectedSize > 0 {
				assert.Equal(t, tt.expectedSize, n)
			}

			result, err := cc.Decode(buf)
			assert.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, result)
			}

			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestCustomCodec_MissingAttribute(t *testing.T) {
	cc := NewCustomCodec(false)
	buf := &bytes.Buffer{}

	cc.AddField("name", codec.Wrap(complex.NewStringCodec()))
	cc.AddField("age", codec.Wrap(&primitive.IntCodec{}))

	data := map[string]any{
		"name": "Alice",
	}

	_, err := cc.Encode(data, buf)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing attribute")
}

func TestCustomCodec_Nested(t *testing.T) {
	buf := &bytes.Buffer{}

	posCodec := NewCustomCodec(true)
	posCodec.AddField("x", codec.Wrap(&primitive.FloatCodec{}))
	posCodec.AddField("y", codec.Wrap(&primitive.FloatCodec{}))

	maincc := NewCustomCodec(true)
	maincc.AddField("id", codec.Wrap(&primitive.IntCodec{}))
	maincc.AddField("pos", codec.Wrap(posCodec))
	maincc.AddField("active", codec.Wrap(&primitive.BoolCodec{}))

	data := map[string]any{
		"id": int32(42),
		"pos": map[string]any{
			"x": float32(1.5),
			"y": float32(2.5),
		},
		"active": true,
	}

	n, err := maincc.Encode(data, buf)

	assert.NoError(t, err)
	assert.True(t, n > 0)

	result, err := maincc.Decode(buf)
	assert.NoError(t, err)

	assert.Equal(t, int32(42), result["id"])

	pos := result["pos"].(map[string]any)
	assert.InDelta(t, float32(1.5), pos["x"].(float32), 0.0001)
	assert.InDelta(t, float32(2.5), pos["y"].(float32), 0.0001)

	assert.Equal(t, true, result["active"])
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_Collections(t *testing.T) {
	cc := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	cc.AddField("strings", codec.Wrap(complex.NewVectorStringCodec()))
	cc.AddField("shorts", codec.Wrap(complex.NewVectorShortCodec()))
	cc.AddField("vectors", codec.Wrap(complex.NewVectorVector3DCodec()))

	data := map[string]any{
		"strings": []string{"a", "b", "c"},
		"shorts":  []int16{1, 2, 3, 4},
		"vectors": []map[string]float32{
			{"x": float32(1), "y": float32(2), "z": float32(3)},
			{"x": float32(4), "y": float32(5), "z": float32(6)},
		},
	}

	n, err := cc.Encode(data, buf)

	assert.NoError(t, err)
	assert.True(t, n > 0)

	result, err := cc.Decode(buf)
	assert.NoError(t, err)

	assert.Equal(t, []string{"a", "b", "c"}, result["strings"])
	assert.Equal(t, []int16{1, 2, 3, 4}, result["shorts"])

	vectors := result["vectors"].([]map[string]float32)
	assert.Len(t, vectors, 2)

	assert.InDelta(t, float32(1), vectors[0]["x"], 0.0001)
	assert.InDelta(t, float32(2), vectors[0]["y"], 0.0001)
	assert.InDelta(t, float32(3), vectors[0]["z"], 0.0001)

	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_ErrorHandling(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "decode empty buffer",
			test: func(t *testing.T) {
				cc := NewCustomCodec(false)
				cc.AddField("test", codec.Wrap(&primitive.IntCodec{}))

				_, err := cc.Decode(&bytes.Buffer{})
				assert.Error(t, err)
			},
		},
		{
			name: "decode empty buffer with boolshortern",
			test: func(t *testing.T) {
				cc := NewCustomCodec(true)
				cc.AddField("test", codec.Wrap(&primitive.IntCodec{}))

				_, err := cc.Decode(&bytes.Buffer{})
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to decode empty flag")
			},
		},
		{
			name: "encode wrong type panic",
			test: func(t *testing.T) {
				cc := NewCustomCodec(false)
				buf := &bytes.Buffer{}
				cc.AddField("age", codec.Wrap(&primitive.IntCodec{}))

				data := map[string]any{
					"age": "not an int",
				}

				assert.Panics(t, func() {
					cc.Encode(data, buf)
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
