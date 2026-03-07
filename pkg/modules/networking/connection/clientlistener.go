package connection

import (
	"net"
)

// ClientListener listens to the ClientConnection
type ClientListener struct {
	listener net.Listener
}

func NewClientListener(addr string) (*ClientListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClientListener{listener: l}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *ClientListener) Accept() (*ClientConnection, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	return NewClientConnection(conn), nil
}

// Close closes the listener. Any blocked Accept operations will be unblocked and return errors.
func (l *ClientListener) Close() error {
	return l.listener.Close()
}
