package packets

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/internal/protection"
)

// Packet header length in bytes.
const HeaderLength = 8

// Represents packets that are sent and received by the server.
type Packet interface {
	ID() int64
	Codecs() []codec.Codec
	Attributes() []string
	ShouldLog() bool
	Unwrap(packetData *bytes.Buffer) (map[string]any, error)
	Wrap(protection *protection.Protection, s2cProxy bool) (*bytes.Buffer, error)
}

// Base packet for concrete packets.
// Provides all necessary functionality and should be embedded to concrete packet implementations.
//
// Concrete packet implementations should only
// provide id, description, codecs, attributes and shouldLog in their constructors.
type BasePacket struct {
	id         int32
	codecs     []codec.Codec
	attributes []string
	shouldLog  bool

	objects []any
	object  map[string]any
}

func NewBasePacket(id int32, codecs []codec.Codec, attributes []string, shouldLog bool) *BasePacket {
	if len(codecs) != len(attributes) {
		panic(fmt.Sprintf(
			"NewBasePacket: codecs and attributes length must be equal; codecs has %d elements, attributes has %d elements",
			len(codecs), len(attributes)))
	}

	attrs := make([]string, len(attributes))
	cdcs := make([]codec.Codec, len(codecs))
	copy(attrs, attributes)
	copy(cdcs, codecs)

	return &BasePacket{
		id:         id,
		codecs:     cdcs,
		attributes: attrs,
		shouldLog:  shouldLog,
		objects:    make([]any, 0),
		object:     make(map[string]any),
	}
}

// Decodes the binary data into individual objects.
func (bp *BasePacket) Unwrap(packetData *bytes.Buffer) (map[string]any, error) {
	for _, c := range bp.codecs {
		decoded, err := c.Decode(packetData)
		if err != nil {
			return nil, fmt.Errorf("BasePacket.Unwrap: failed to unwrap: %w", err)
		}
		bp.objects = append(bp.objects, decoded)
	}

	return bp.implement(), nil
}

// Encodes all the objects into binary data for the packet payload.
func (bp *BasePacket) Wrap(protection protection.Protection) (*bytes.Buffer, error) {
	if protection == nil {
		panic("BasePacket.Wrap: nil protection is passed")
	}

	packetData := &bytes.Buffer{}
	dataLen := HeaderLength

	for i, c := range bp.codecs {
		n, err := c.Encode(bp.objects[i], packetData)
		if err != nil {
			return nil, fmt.Errorf("BasePacket.Wrap: failed to encode packed data: %w", err)
		}
		dataLen += n
	}

	intCodec := primitive.IntCodec{}
	encryptedData := protection.Encrypt(packetData.Bytes())

	packetData = &bytes.Buffer{}
	_, err := intCodec.Encode(int32(dataLen), packetData)
	if err != nil {
		return nil, fmt.Errorf("BasePacket.Wrap: failed to encode data length: %w", err)
	}

	_, err = intCodec.Encode(bp.id, packetData)
	if err != nil {
		return nil, fmt.Errorf("BasePacket.Wrap: failed to encode packet id: %w", err)
	}

	_, err = packetData.Write(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("BasePacket.Wrap: failed to write encoded packet data: %w", err)
	}

	return packetData, nil
}

// Implements the packet object based on the attribute key list and the decoded object list.
func (bp *BasePacket) implement() map[string]any {
	clear(bp.object)
	for i, obj := range bp.objects {
		bp.object[bp.attributes[i]] = obj
	}
	return bp.object
}
