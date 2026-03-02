package connection

import "time"

// Connection denotes all network connections.
type Connection interface {
	// Read reads exactly n bytes
	Read(n int) ([]byte, error)

	// Write writes data to connection
	Write(data []byte) (int, error)

	// Close closes the connection
	Close() error

	// SetDeadline sets read/write deadlines
	SetDeadline(t time.Time) error

	// SetReadDeadline sets read deadline
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets write deadline
	SetWriteDeadline(t time.Time) error
}
