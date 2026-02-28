package networking

import (
	"bytes"
	"compress/flate"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

const (
	DefaultSocketRetryDelay = 2 * time.Second
	DefaultSocketMaxRetries = 3
	packetHeaderLength      = 8
)

// TankiSocket handles asynchronous network communication with the game server
type TankiSocket struct {
	// endpoint is the target game server address (host:port)
	endpoint *Address

	// proxy is an optional SOCKS5 proxy address (host:port with optional auth)
	proxy *Address

	// protection provides encryption/decryption for packet data
	protection protection.Protection

	// emergencyHalt is a context that signals all goroutines to stop
	emergencyHalt context.Context

	// onDataReceived is called when a new packet is successfully received and parsed
	onDataReceived func(packet packets.Packet) error

	// onSocketClose is called when the socket is closed (due to error or intentionally)
	onSocketClose func(err error, source string, details string)

	// mu protects concurrent access to conn
	mu sync.RWMutex

	// conn is the current active network connection
	conn net.Conn

	// cancel cancels the emergencyHalt context
	cancel context.CancelFunc

	// processingDone is closed when ProcessSocket goroutine exits
	processingDone chan struct{}

	// socketRetryDelay is the base delay between connection attempts (increases exponentially)
	socketRetryDelay time.Duration

	// socketMaxRetries is the maximum number of connection attempts
	socketMaxRetries int
}

// NewTankiSocket creates a new socket instance.
//
// Parameters:
//   - endpoint: target game server address
//   - protection: encryption/decryption object
//   - proxy: optional SOCKS5 proxy address (nil for direct connection)
//   - ctx: parent context for cancellation propagation
//   - onDataReceived: callback for incoming packets
//   - onSocketClose: callback for socket closure events
//
// The socket starts in disconnected state. Call ProcessSocket() to begin operation.
func NewTankiSocket(
	endpoint *Address,
	protection protection.Protection,
	proxy *Address,
	ctx context.Context,
	onDataReceived func(packet packets.Packet) error,
	onSocketClose func(err error, source string, details string),
) *TankiSocket {
	ctx, cancel := context.WithCancel(ctx)
	s := &TankiSocket{
		endpoint:         endpoint,
		protection:       protection,
		proxy:            proxy,
		emergencyHalt:    ctx,
		onDataReceived:   onDataReceived,
		onSocketClose:    onSocketClose,
		cancel:           cancel,
		processingDone:   make(chan struct{}),
		socketRetryDelay: DefaultSocketRetryDelay,
		socketMaxRetries: DefaultSocketMaxRetries,
	}

	return s
}

// SetSocketRetryDelay sets socket retry delay.
// DefaultSocketRetryDelay is used by default.
func (s *TankiSocket) SetSocketRetryDelay(delay time.Duration) {
	s.socketRetryDelay = delay
}

// SetSocketMaxRetries sets socket max retries.
// DefaultSocketMaxRetries is used by default.
func (s *TankiSocket) SetSocketMaxRetries(maxRetries int) {
	s.socketMaxRetries = maxRetries
}

// readFull reads exactly n bytes from the connection.
// It assumes the caller has already ensured that s.conn is not nil.
func (s *TankiSocket) readFull(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(s.conn, buf)
	return buf, err
}

// connect establishes connection to endpoint with retry and backoff.
// Returns nil on success, otherwise an error.
func (s *TankiSocket) connect() error {
	var lastErr error

	for attempt := 0; attempt < s.socketMaxRetries; attempt++ {
		select {
		case <-s.emergencyHalt.Done():
			return context.Canceled
		default:
		}

		s.closeConn()

		var conn net.Conn
		var err error

		targetAddr := fmt.Sprintf("%s:%d", s.endpoint.Host, s.endpoint.Port)

		if s.proxy != nil {
			conn, err = s.connectViaProxy(targetAddr)
		} else {
			conn, err = s.connectDirect(targetAddr)
		}

		if err == nil {
			// Enable TCP_NODELAY for low latency
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.SetNoDelay(true)
			}

			s.mu.Lock()
			s.conn = conn
			s.mu.Unlock()

			return nil
		}

		lastErr = err

		if attempt < s.socketMaxRetries-1 {
			backoff := s.socketRetryDelay * time.Duration(1<<uint(attempt))
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", s.socketMaxRetries, lastErr)
}

// connectDirect establishes a direct TCP or TLS connection.
func (s *TankiSocket) connectDirect(targetAddr string) (net.Conn, error) {
	if s.endpoint.Port == 443 {
		return tls.Dial("tcp", targetAddr, &tls.Config{
			ServerName: s.endpoint.Host,
		})
	}
	return net.DialTimeout("tcp", targetAddr, 10*time.Second)
}

// connectViaProxy establishes a connection through a SOCKS5 proxy.
// If the target port is 443, the connection is wrapped with TLS after the proxy tunnel.
func (s *TankiSocket) connectViaProxy(targetAddr string) (net.Conn, error) {
	proxyAddr := fmt.Sprintf("%s:%d", s.proxy.Host, s.proxy.Port)
	auth := &proxy.Auth{
		User:     s.proxy.Username,
		Password: s.proxy.Password,
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	conn, err := dialer.Dial("tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("SOCKS5 dial failed: %w", err)
	}

	// If the target is TLS, upgrade the connection
	if s.endpoint.Port == 443 {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName: s.endpoint.Host,
		})
		// Perform the TLS handshake explicitly
		if err := tlsConn.Handshake(); err != nil {
			conn.Close()
			return nil, fmt.Errorf("TLS handshake failed: %w", err)
		}
		return tlsConn, nil
	}

	return conn, nil
}

// ProcessSocket is the main socket processing loop.
func (s *TankiSocket) ProcessSocket() {
	defer close(s.processingDone)

	if err := s.connect(); err != nil {
		s.onSocketClose(err, "TankiSocket.connect", "Failed to connect")
		return
	}

	for {
		select {
		case <-s.emergencyHalt.Done():
			return
		default:
			packetLen, packetID, isCompressed, err := s.readPacketHeader()
			if err != nil {
				s.onSocketClose(err, "TankiSocket.processSocket",
					"Connected | Failed to read header")
				return
			}

			packetDataLen := packetLen - packetHeaderLength
			var encryptedData []byte

			if packetDataLen > 0 {
				encryptedData, err = s.readFull(int(packetDataLen))
				if err != nil {
					s.onSocketClose(err, "TankiSocket.processSocket",
						fmt.Sprintf("Connected | Packet Length: %d | Failed to read data", packetLen))
					return
				}
			}

			if err := s.processPacket(packetID, encryptedData, isCompressed); err != nil {
				s.onSocketClose(err, "TankiSocket.processSocket",
					fmt.Sprintf("Connected | Packet ID: %d | Failed to process", packetID))
				return
			}
		}
	}
}

// readPacketHeader reads the 8-byte header (4 bytes length+flags, 4 bytes packet ID).
func (s *TankiSocket) readPacketHeader() (int32, int32, bool, error) {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()

	if conn == nil {
		return 0, 0, false, fmt.Errorf("not connected")
	}

	// Read first 4 bytes: length and compression flag
	headerBytes, err := s.readFull(4)
	if err != nil {
		return 0, 0, false, err
	}

	header := binary.BigEndian.Uint32(headerBytes)
	isCompressed := (header>>24)&0x40 != 0
	packetLen := int32(header & 0x00FFFFFF)

	// Read next 4 bytes: packet ID
	idBytes, err := s.readFull(4)
	if err != nil {
		return 0, 0, false, err
	}

	packetID := int32(binary.BigEndian.Uint32(idBytes))

	return packetLen, packetID, isCompressed, nil
}

// processPacket decrypts, decompresses (if needed), and delivers the packet.
func (s *TankiSocket) processPacket(packetID int32, encryptedData []byte, isCompressed bool) error {
	decrypted := s.protection.Decrypt(encryptedData)

	if isCompressed {
		r := flate.NewReader(bytes.NewReader(decrypted))
		defer r.Close()

		var err error
		decrypted, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("decompression failed: %w", err)
		}
	}

	packet, err := s.fitPacket(packetID, decrypted)
	if err != nil {
		return err
	}

	return s.onDataReceived(packet)
}

// fitPacket converts raw data (decrypted) into a packet object using the packet registry.
func (s *TankiSocket) fitPacket(packetID int32, data []byte) (packets.Packet, error) {
	packet := packets.Get(packetID)
	if packet == nil {
		return packets.NewUnknownPacket(packetID, data), nil
	}

	if _, err := packet.Unwrap(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	return packet, nil
}

// Send writes raw data to the connection.
func (s *TankiSocket) Send(data []byte) error {
	select {
	case <-s.emergencyHalt.Done():
		return context.Canceled
	default:
	}

	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	_, err := conn.Write(data)
	return err
}

// SendBatch concatenates multiple data chunks and sends them in one write.
func (s *TankiSocket) SendBatch(packetsData [][]byte) error {
	var batch []byte
	for _, p := range packetsData {
		batch = append(batch, p...)
	}
	return s.Send(batch)
}

// closeConn safely closes and removes the current connection.
func (s *TankiSocket) closeConn() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

// Close gracefully shuts down the socket and waits for the processing loop to exit.
func (s *TankiSocket) Close() error {
	s.cancel()
	s.closeConn()

	<-s.processingDone
	return nil
}
