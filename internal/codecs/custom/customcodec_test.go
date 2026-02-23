package custom

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/internal/codecs"
	"github.com/the-new-day/probogo/internal/codecs/complex"
	"github.com/the-new-day/probogo/internal/codecs/primitive"
)

func TestCustomCodec_WithoutBoolshortern(t *testing.T) {
	codec := NewCustomCodec(false)
	buf := &bytes.Buffer{}

	// Примитивные без конструктора, сложные — через конструкторы
	codec.AddField("name", codecs.Wrap(complex.NewStringCodec()))
	codec.AddField("age", codecs.Wrap(&primitive.IntCodec{}))
	codec.AddField("active", codecs.Wrap(&primitive.BoolCodec{}))

	data := map[string]any{
		"name":   "Alice",
		"age":    int32(30),
		"active": true,
	}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	expectedSize := 1 + 4 + len("Alice") + 4 + 1
	assert.Equal(t, expectedSize, n)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["name"], result["name"])
	assert.Equal(t, data["age"], result["age"])
	assert.Equal(t, data["active"], result["active"])
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_WithBoolshortern_NotEmpty(t *testing.T) {
	codec := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	codec.AddField("x", codecs.Wrap(&primitive.FloatCodec{}))
	codec.AddField("y", codecs.Wrap(&primitive.FloatCodec{}))

	data := map[string]any{
		"x": float32(1.5),
		"y": float32(2.7),
	}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1+4+4, n)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.InDelta(t, float32(1.5), result["x"].(float32), 0.0001)
	assert.InDelta(t, float32(2.7), result["y"].(float32), 0.0001)
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_WithBoolshortern_Empty(t *testing.T) {
	codec := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	codec.AddField("x", codecs.Wrap(&primitive.FloatCodec{}))
	codec.AddField("y", codecs.Wrap(&primitive.FloatCodec{}))

	data := map[string]any{}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_MissingAttribute(t *testing.T) {
	codec := NewCustomCodec(false)
	buf := &bytes.Buffer{}

	codec.AddField("name", codecs.Wrap(complex.NewStringCodec()))
	codec.AddField("age", codecs.Wrap(&primitive.IntCodec{}))

	data := map[string]any{
		"name": "Alice",
	}

	_, err := codec.Encode(data, buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing attribute")
}

func TestCustomCodec_MultipleFields(t *testing.T) {
	codec := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	// Микс: примитивы без конструктора, сложные — с конструкторами
	codec.AddField("id", codecs.Wrap(&primitive.LongCodec{}))
	codec.AddField("name", codecs.Wrap(complex.NewStringCodec()))
	codec.AddField("position", codecs.Wrap(complex.NewVector3DCodec()))
	codec.AddField("tags", codecs.Wrap(complex.NewVectorStringCodec()))
	codec.AddField("health", codecs.Wrap(&primitive.ShortCodec{}))

	data := map[string]any{
		"id":   int64(12345),
		"name": "EnemyTank",
		"position": map[string]any{
			"x": float32(10.5),
			"y": float32(20.3),
			"z": float32(5.0),
		},
		"tags":   []string{"boss", "elite", "armored"},
		"health": int16(100),
	}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)

	assert.Equal(t, int64(12345), result["id"])
	assert.Equal(t, "EnemyTank", result["name"])

	pos := result["position"].(map[string]any)
	assert.InDelta(t, float32(10.5), pos["x"].(float32), 0.0001)
	assert.InDelta(t, float32(20.3), pos["y"].(float32), 0.0001)
	assert.InDelta(t, float32(5.0), pos["z"].(float32), 0.0001)

	tags := result["tags"].([]string)
	assert.Equal(t, "boss", tags[0])
	assert.Equal(t, "elite", tags[1])
	assert.Equal(t, "armored", tags[2])

	assert.Equal(t, int16(100), result["health"])
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_WithNestedCustomCodecs(t *testing.T) {
	buf := &bytes.Buffer{}

	// Вложенный кастомный кодек
	posCodec := NewCustomCodec(true)
	posCodec.AddField("x", codecs.Wrap(&primitive.FloatCodec{}))
	posCodec.AddField("y", codecs.Wrap(&primitive.FloatCodec{}))

	// Главный кодек
	mainCodec := NewCustomCodec(true)
	mainCodec.AddField("id", codecs.Wrap(&primitive.IntCodec{}))
	mainCodec.AddField("pos", codecs.Wrap(posCodec))
	mainCodec.AddField("active", codecs.Wrap(&primitive.BoolCodec{}))

	data := map[string]any{
		"id": int32(42),
		"pos": map[string]any{
			"x": float32(1.5),
			"y": float32(2.5),
		},
		"active": true,
	}

	n, err := mainCodec.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	result, err := mainCodec.Decode(buf)
	assert.NoError(t, err)

	assert.Equal(t, int32(42), result["id"])

	pos := result["pos"].(map[string]any)
	assert.InDelta(t, float32(1.5), pos["x"].(float32), 0.0001)
	assert.InDelta(t, float32(2.5), pos["y"].(float32), 0.0001)

	assert.Equal(t, true, result["active"])
	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_MultipleCodecs(t *testing.T) {
	codec := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	codec.AddField("strings", codecs.Wrap(complex.NewVectorStringCodec()))
	codec.AddField("shorts", codecs.Wrap(complex.NewVectorShortCodec()))
	codec.AddField("vectors", codecs.Wrap(complex.NewVectorVector3DCodec()))

	data := map[string]any{
		"strings": []string{"a", "b", "c"},
		"shorts":  []int16{1, 2, 3, 4},
		"vectors": []map[string]any{
			{"x": float32(1), "y": float32(2), "z": float32(3)},
			{"x": float32(4), "y": float32(5), "z": float32(6)},
		},
	}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)

	strings := result["strings"].([]string)
	assert.Equal(t, "a", strings[0])
	assert.Equal(t, "b", strings[1])
	assert.Equal(t, "c", strings[2])

	shorts := result["shorts"].([]int16)
	assert.Equal(t, int16(1), shorts[0])
	assert.Equal(t, int16(2), shorts[1])
	assert.Equal(t, int16(3), shorts[2])
	assert.Equal(t, int16(4), shorts[3])

	vectors := result["vectors"].([]map[string]any)
	assert.Len(t, vectors, 2)

	v1 := vectors[0]
	assert.InDelta(t, float32(1), v1["x"].(float32), 0.0001)
	assert.InDelta(t, float32(2), v1["y"].(float32), 0.0001)
	assert.InDelta(t, float32(3), v1["z"].(float32), 0.0001)

	assert.Equal(t, 0, buf.Len())
}

func TestCustomCodec_ErrorHandling(t *testing.T) {
	t.Run("decode with empty buffer", func(t *testing.T) {
		codec := NewCustomCodec(false)
		buf := &bytes.Buffer{}
		codec.AddField("test", codecs.Wrap(&primitive.IntCodec{}))

		_, err := codec.Decode(buf)
		assert.Error(t, err)
	})

	t.Run("decode with boolshortern and empty buffer", func(t *testing.T) {
		codec := NewCustomCodec(true)
		buf := &bytes.Buffer{}
		codec.AddField("test", codecs.Wrap(&primitive.IntCodec{}))

		_, err := codec.Decode(buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode empty flag")
	})

	t.Run("encode with wrong type", func(t *testing.T) {
		codec := NewCustomCodec(false)
		buf := &bytes.Buffer{}
		codec.AddField("age", codecs.Wrap(&primitive.IntCodec{}))

		data := map[string]any{
			"age": "not an int",
		}

		assert.Panics(t, func() {
			codec.Encode(data, buf)
		})
	})
}

func TestCustomCodec_NoFields(t *testing.T) {
	codec := NewCustomCodec(true)
	buf := &bytes.Buffer{}

	data := map[string]any{"unused": "value"}

	n, err := codec.Encode(data, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
