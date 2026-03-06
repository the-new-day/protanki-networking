// Package networking provides components for network communication with the game server.
package networking

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

// PacketStream transforms raw network data into game packets and vice versa.
// It handles packet framing, encryption/decryption, compression/decompression,
// and packet type resolution through a registry.
type PacketStream struct {
	conn           connection.Connection
	protection     protection.Protection
	packetRegistry *packets.PacketRegistry

	onActivateProtection func([]byte)

	mu sync.RWMutex
}

// NewPacketStream creates a new PacketStream that reads from and writes to the given connection.
// The stream uses the provided protection for encryption/decryption and the registry
// to resolve packet types by their IDs.
func NewPacketStream(
	conn connection.Connection,
	protection protection.Protection,
	packetRegistry *packets.PacketRegistry,
	onActivateProtection func([]byte),
) *PacketStream {
	return &PacketStream{
		conn:                 conn,
		protection:           protection,
		packetRegistry:       packetRegistry,
		onActivateProtection: onActivateProtection,
		mu:                   sync.RWMutex{},
	}
}

// PacketResult represents either a successfully parsed packet or an error that occurred
// during packet processing. It is used by the Packets method to deliver results
// asynchronously.
type PacketResult struct {
	// Packet is the successfully parsed packet. Nil if Err is not nil.
	Packet packets.Packet

	// Err describes an error that occurred while reading or processing the packet.
	// If non-nil, the Packet field should be ignored.
	Err error

	// ID is the packet ID read from the connection.
	ID int32

	// Length is the packet length read from the connection (without compression bit).
	Length int32

	// Data is the raw bytes representing the packet data (without ID and length) (decrypted).
	Data []byte

	// RawHex is the raw bytes received from the connection (not decrypted, not decompressed).
	// Includes length and ID.
	RawHex []byte

	WasCompressed bool
}

// Send encodes and sends a packet through the underlying connection.
// It uses the stream's protection to encrypt the packet data before transmission.
// Returns an error if encoding fails or the write operation fails.
func (ps *PacketStream) Send(packet packets.Packet) error {
	var payload *bytes.Buffer
	var err error

	ps.mu.Lock()
	payload, err = packet.Wrap(ps.protection)
	ps.mu.Unlock()

	if err != nil {
		return err
	}

	_, err = ps.conn.Write(payload.Bytes())
	return err
}

// SendRaw sends raw data without encryption.
func (ps *PacketStream) SendRaw(data []byte) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, err := ps.conn.Write(data)
	return err
}

// SendRawEncrypted encrypts and sends raw data.
func (ps *PacketStream) SendRawEncrypted(data []byte) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	data = ps.protection.Encrypt(data)
	_, err := ps.conn.Write(data)
	return err
}

// Inbound returns a channel that delivers parsed packets as they arrive.
// The method runs in a separate goroutine and continues until:
//   - the context is cancelled
//   - a fatal read error occurs (connection closed, etc.)
//   - the underlying connection returns an unrecoverable error
//
// The channel is closed when the stream stops.
//
// Example usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	for res := range stream.Inbound(ctx) {
//		if res.Err != nil {
//			log.Printf("packet error: %v", res.Err)
//			continue
//		}
//		handlePacket(res.Packet)
//	}
func (ps *PacketStream) Inbound(ctx context.Context) <-chan PacketResult {
	ch := make(chan PacketResult)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			result := PacketResult{}

			rawHex, packetLen, packetID, isCompressed, err := ps.readPacketHeader()
			result.Length = packetLen
			result.ID = packetID
			result.RawHex = rawHex
			result.WasCompressed = isCompressed

			if err != nil {
				result.Err = err
				ch <- result
				return
			}

			packetDataLen := packetLen - packets.HeaderLength
			var encryptedData []byte

			if packetDataLen > 0 {
				encryptedData, err = ps.conn.Read(int(packetDataLen))
				if err != nil {
					result.Err = err
					ch <- result
					return
				}
				result.RawHex = append(result.RawHex, encryptedData...)
			}

			ps.mu.Lock()
			decryptedData := ps.protection.Decrypt(encryptedData)
			ps.mu.Unlock()

			result.Data = decryptedData

			packet, err := ps.processPacket(packetID, decryptedData, isCompressed)

			if err != nil {
				result.Err = err
				ch <- result
				continue
			}

			if isCompressed {
				packet.SetCompress(true)
			}

			if packetID == packets.ActivateProtectionID {
				keys := packets.Attr[[]byte]("keys", packet)
				ps.onActivateProtection(keys)
			}

			result.Packet = packet
			ch <- result
		}
	}()

	return ch
}

// readPacketHeader reads and parses the 8-byte packet header.
// Header format:
//   - bytes 0-3: packet length with compression flag in bit 24
//   - bytes 4-7: packet ID
//
// Returns:
//   - rawHex: raw bytes for the header
//   - packetLen: total packet length including header
//   - packetID: unique packet identifier
//   - isCompressed: whether the payload is compressed with zlib
//   - error: any read or parsing error
func (ps *PacketStream) readPacketHeader() ([]byte, int32, int32, bool, error) {
	rawHex := []byte{}

	if ps.conn == nil {
		return rawHex, 0, 0, false, connection.ErrNotConnected
	}

	// read first 4 bytes: length and compression flag
	lengthBytes, err := ps.conn.Read(4)
	if err != nil {
		return rawHex, 0, 0, false, err
	}

	length := binary.BigEndian.Uint32(lengthBytes)
	isCompressed := (length>>24)&0x40 != 0
	packetLen := int32(length & 0x00FFFFFF) // nullifies the compression bit

	// read next 4 bytes: packet ID
	idBytes, err := ps.conn.Read(4)
	if err != nil {
		return rawHex, 0, 0, false, err
	}

	packetID := int32(binary.BigEndian.Uint32(idBytes))

	rawHex = append(rawHex, lengthBytes...)
	rawHex = append(rawHex, idBytes...)
	return rawHex, packetLen, packetID, isCompressed, nil
}

// processPacket decompresses (if needed) and creates a packet object.
// Returns the parsed packet or an error if any processing step fails.
func (ps *PacketStream) processPacket(packetID int32, decryptedData []byte, isCompressed bool) (packets.Packet, error) {
	if isCompressed {
		var err error
		decryptedData, err = packets.Decompress(decryptedData)
		if err != nil {
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
	}

	packet, err := ps.fitPacket(packetID, decryptedData)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

// fitPacket converts raw decrypted data into a packet object using the packet registry.
// If no packet type is registered for the given ID, returns an UnknownPacket containing
// the raw data.
func (ps *PacketStream) fitPacket(packetID int32, data []byte) (packets.Packet, error) {
	packet := ps.packetRegistry.Get(packetID)
	if packet == nil {
		return packets.NewUnknownPacket(packetID, data), nil
	}

	if _, err := packet.Unwrap(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	return packet, nil
}

func (ps *PacketStream) ActivateProtection(keys []byte) {
	ps.protection.Activate(keys)
}
