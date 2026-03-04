package packets

import (
	"bytes"
	"fmt"

	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/modules/protection"
)

// Packet header length in bytes.
const HeaderLength = 8

// Packet represents packets that are sent and received by the server.
type Packet interface {
	// ID returns packet ID, usually set during initialization.
	ID() int32

	// Len returns packet length in bytes, including header.
	// Fills during Unwrap, bufore that it can be zero.
	Len() int

	// Unwrap decodes the binary data into individual objects.
	// May store decoded objects for future use.
	Unwrap(packetData *bytes.Buffer) (map[string]any, error)

	// Wrap encodes all the objects into binary data for the packet payload.
	// Does not affect inner state of the packet, but may affect inner state of Protection.
	Wrap(protection protection.Protection) (*bytes.Buffer, error)

	// Attr returns value of the given attrubute or panics if it doesn't exists.
	// It also works for attributes added by Set().
	Attr(name string) any

	// Set sets value for the attribute (it's possible to add new attribute, but it won't be wrapped).
	// It does not perform type assertions, encryption/decryption etc.
	Set(name string, value any)

	// EncryptedRawData returns raw bytes that contain (probably encrypted) with packet data from the server.
	EncryptedRawData() []byte

	SetEncryptedRawData([]byte)

	// RawData returns packet representation in bytes (decrypted).
	// Fills during Unwrap, before that can be empty.
	RawData() []byte
}

// Base packet for concrete packets.
// Provides all necessary functionality and should be embedded to concrete packet implementations.
//
// Concrete packet implementations should only
// provide id, codecs, attributes in their constructors
// unless they need some specific implementation, as in the example:
//
//	type FireEndOutPacket struct {
//		packets.BasePacket
//	}
//
//	func NewFireEndOutPacket() *FireEndOutPacket {
//		codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{})}
//		attributes := []string{"clientTime"}
//		var id int32 = packets.FireEndOutID
//
//		return &FireEndOutPacket{
//			BasePacket: *packets.NewBasePacket(id, codecs, attributes),
//		}
//	}
type BasePacket struct {
	id         int32
	codecs     []codec.Codec
	attributes []string

	encryptedRawData []byte
	rawData          []byte

	objects []any

	object map[string]any
}

func NewBasePacket(id int32, codecs []codec.Codec, attributes []string) *BasePacket {
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
		objects:    make([]any, 0),
		object:     make(map[string]any),
	}
}

func (bp *BasePacket) Unwrap(packetData *bytes.Buffer) (map[string]any, error) {
	buf := make([]byte, packetData.Len())
	copy(buf, packetData.Bytes())

	for _, c := range bp.codecs {
		decoded, err := c.Decode(packetData)
		if err != nil {
			return nil, fmt.Errorf("BasePacket.Unwrap: packet ID: %d | failed to unwrap: %w", bp.id, err)
		}
		bp.objects = append(bp.objects, decoded)
	}

	bp.rawData = buf
	return bp.populate(), nil
}

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

	encryptedData := protection.Encrypt(packetData.Bytes())

	packetData = &bytes.Buffer{}
	intCodec := primitive.IntCodec{}

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

func (bp *BasePacket) Attr(name string) any {
	bp.populate()
	value, ok := bp.object[name]
	if !ok {
		panic(fmt.Sprintf("BasePacket.Get: attribute %q not found. ID: %d", name, bp.ID()))
	}
	return value
}

func (bp *BasePacket) Set(name string, value any) {
	bp.object[name] = value
}

// populate fills the packet object based on the attribute key list and the decoded object list.
func (bp *BasePacket) populate() map[string]any {
	for i, obj := range bp.objects {
		bp.object[bp.attributes[i]] = obj
	}
	return bp.object
}

func (bp *BasePacket) ID() int32 {
	return bp.id
}

func (bp *BasePacket) EncryptedRawData() []byte {
	return bp.rawData
}

func (bp *BasePacket) RawData() []byte {
	return bp.rawData
}

func (bp *BasePacket) SetEncryptedRawData(data []byte) {
	buf := make([]byte, len(data))
	copy(buf, data)
	bp.rawData = buf
}

func (bp *BasePacket) Len() int {
	return len(bp.rawData)
}
