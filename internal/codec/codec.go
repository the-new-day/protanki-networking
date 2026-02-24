package codec

import (
	"bytes"
	"fmt"
)

// Interface for all codecs, whether primitive or complex.
type Codec interface {
	Decode(buf *bytes.Buffer) (any, error)
	Encode(rawValue any, buf *bytes.Buffer) (int, error)
}

// Codec for concrete type.
type TypedCodec[T any] interface {
	Decode(buf *bytes.Buffer) (T, error)
	Encode(value T, buf *bytes.Buffer) (int, error)
}

// Proxy for typed codecs: Encode() checks that passed value is of type T.
type CodecWrapper[T any] struct {
	typedCodec TypedCodec[T]
}

// Decodes value from reader using underlying TypedCodec.
// This method simply transfers the call to TypedCodec.
func (bc *CodecWrapper[T]) Decode(buf *bytes.Buffer) (any, error) {
	return bc.typedCodec.Decode(buf)
}

// Encodes value using underlying TypedCodec.
// Checks if rawValue is of type T, if it's not, panics.
// Returns number of bytes written and an error.
func (bc *CodecWrapper[T]) Encode(rawValue any, buf *bytes.Buffer) (int, error) {
	value, ok := rawValue.(T)
	if !ok {
		var zero T
		panic(fmt.Errorf("CodecWrapper: expected type %T, got %T", zero, rawValue))
	}
	return bc.typedCodec.Encode(value, buf)
}

// Creates an instance of CodecWrapper[T] for specified TypedCodec[T].
func Wrap[T any](tc TypedCodec[T]) *CodecWrapper[T] {
	return &CodecWrapper[T]{typedCodec: tc}
}
