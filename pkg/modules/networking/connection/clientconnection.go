package connection

import (
	"io"
	"net"
	"time"
)

// ClientConnection is an adapter for Connection interface.
// It holds established net.Conn object and transfers calls to it.
// It's used to handler connection with the game client.
type ClientConnection struct {
	conn net.Conn
}

func NewClientConnection(conn net.Conn) *ClientConnection {
	return &ClientConnection{conn}
}

func (c *ClientConnection) Read(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(c.conn, buf)
	return buf, err
}

func (c *ClientConnection) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *ClientConnection) Close() error {
	return c.conn.Close()
}

func (c *ClientConnection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *ClientConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *ClientConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
