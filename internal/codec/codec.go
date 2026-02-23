package codec

import (
	"fmt"
	"io"
)

// Interface for all codecs, whether primitive or complex.
type Codec interface {
	Decode(reader io.Reader) (any, error)
	Encode(rawValue any, writer io.Writer) (int, error)
}

// Codec for concrete type.
type TypedCodec[T any] interface {
	Decode(reader io.Reader) (T, error)
	Encode(value T, writer io.Writer) (int, error)
}

// Proxy for typed codecs: Encode() checks that passed value is of type T.
type CodecWrapper[T any] struct {
	typedCodec TypedCodec[T]
}

func (bc *CodecWrapper[T]) Decode(reader io.Reader) (any, error) {
	return bc.typedCodec.Decode(reader)
}

func (bc *CodecWrapper[T]) Encode(rawValue any, writer io.Writer) (int, error) {
	value, ok := rawValue.(T)
	if !ok {
		var zero T
		panic(fmt.Errorf("CodecWrapper: expected type %T, got %T", zero, rawValue))
	}
	return bc.typedCodec.Encode(value, writer)
}

// Creates an instance of CodecWrapper[T] for specified TypedCodec[T].
func Wrap[T any](tc TypedCodec[T]) *CodecWrapper[T] {
	return &CodecWrapper[T]{typedCodec: tc}
}
