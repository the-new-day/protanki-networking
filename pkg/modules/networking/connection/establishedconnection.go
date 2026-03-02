package connection

import (
	"io"
	"net"
	"time"
)

// EstablishedConnection is an adapter for Connection interface.
// It holds established net.Conn object and transfers calls to it.
type EstablishedConnection struct {
	conn net.Conn
}

func NewEstablishedConnection(conn net.Conn) *EstablishedConnection {
	return &EstablishedConnection{conn}
}

func (c *EstablishedConnection) Read(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(c.conn, buf)
	return buf, err
}

func (c *EstablishedConnection) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *EstablishedConnection) Close() error {
	return c.conn.Close()
}

func (c *EstablishedConnection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *EstablishedConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *EstablishedConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
