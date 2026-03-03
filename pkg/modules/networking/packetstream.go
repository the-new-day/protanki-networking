// Package networking provides components for network communication with the game server.
package networking

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/binary"
	"fmt"
	"io"
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

	mu sync.RWMutex
}

// NewPacketStream creates a new PacketStream that reads from and writes to the given connection.
// The stream uses the provided protection for encryption/decryption and the registry
// to resolve packet types by their IDs.
func NewPacketStream(
	conn connection.Connection,
	protection protection.Protection,
	packetRegistry *packets.PacketRegistry,
) *PacketStream {
	return &PacketStream{conn, protection, packetRegistry, sync.RWMutex{}}
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
}

// Send encodes and sends a packet through the underlying connection.
// It uses the stream's protection to encrypt the packet data before transmission.
// Returns an error if encoding fails or the write operation fails.
func (ps *PacketStream) Send(packet packets.Packet) error {
	ps.mu.Lock()

	payload, err := packet.Wrap(ps.protection)
	if err != nil {
		return err
	}

	ps.mu.Unlock()

	_, err = ps.conn.Write(payload.Bytes())
	return err
}

// Inbound returns a channel that delivers parsed packets as they arrive.
// The method runs in a separate goroutine and continues until:
//   - the context is cancelled
//   - a fatal read error occurs (connection closed, etc.)
//   - the underlying connection returns an unrecoverable error
//
// Non-fatal errors (like decompression failures) are sent as PacketResult with Err set,
// and the stream continues processing subsequent packets.
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

			packetLen, packetID, isCompressed, err := ps.readPacketHeader()

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
			}

			packet, err := ps.processPacket(packetID, encryptedData, isCompressed)
			if err != nil {
				result.Err = err
				ch <- result
				continue
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
//   - packetLen: total packet length including header
//   - packetID: unique packet identifier
//   - isCompressed: whether the payload is compressed with zlib
//   - error: any read or parsing error
func (ps *PacketStream) readPacketHeader() (int32, int32, bool, error) {
	if ps.conn == nil {
		return 0, 0, false, connection.ErrNotConnected
	}

	// read first 4 bytes: length and compression flag
	headerBytes, err := ps.conn.Read(4)
	if err != nil {
		return 0, 0, false, err
	}

	header := binary.BigEndian.Uint32(headerBytes)
	isCompressed := (header>>24)&0x40 != 0
	packetLen := int32(header & 0x00FFFFFF) // nullifies the compression bit

	// read next 4 bytes: packet ID
	idBytes, err := ps.conn.Read(4)
	if err != nil {
		return 0, 0, false, err
	}

	packetID := int32(binary.BigEndian.Uint32(idBytes))

	return packetLen, packetID, isCompressed, nil
}

// processPacket decrypts, decompresses (if needed), and creates a packet object.
// Returns the parsed packet or an error if any processing step fails.
func (ps *PacketStream) processPacket(packetID int32, encryptedData []byte, isCompressed bool) (packets.Packet, error) {
	decrypted := ps.protection.Decrypt(encryptedData)

	if isCompressed {
		r := flate.NewReader(bytes.NewReader(decrypted))
		defer r.Close()

		var err error
		decrypted, err = io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
	}

	packet, err := ps.fitPacket(packetID, decrypted)
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

// ActivateProtection call Activate(keys) on the underlying Protection instance.
func (ps *PacketStream) ActivateProtection(keys []byte) {
	ps.protection.Activate(keys)
}
