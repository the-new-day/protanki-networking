package connection

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/the-new-day/probogo/pkg/modules/networking/dial"
)

const (
	DefaultConnectionRetryDelay = 2 * time.Second
	DefaultConnectionMaxRetries = 3

	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultKeepAlivePeriod = 30 * time.Second
)

var (
	ErrNotConnected = errors.New("networking: not connected")
	ErrShuttingDown = errors.New("networking: connection is shutting down")
)

// GameConnection manages connection with the game server
// and allows sending and receiving data.
type GameConnection struct {
	endpoint string
	dialer   dial.Dialer

	mu   sync.RWMutex
	conn net.Conn

	ctx    context.Context
	cancel context.CancelFunc // controls the main connection context

	retryDelay   time.Duration // base delay between connection attempts (increases exponentially)
	maxRetries   int           // max number of connection attempts
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewGameConnection creates new GameConnection object with unconnected state.
func NewGameConnection(endpoint string, dialer dial.Dialer) *GameConnection {
	ctx, cancel := context.WithCancel(context.Background())

	return &GameConnection{
		endpoint: endpoint,
		dialer:   dialer,

		ctx:    ctx,
		cancel: cancel,

		retryDelay:   DefaultConnectionRetryDelay,
		maxRetries:   DefaultConnectionMaxRetries,
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
	}
}

func (c *GameConnection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil
}

// Read reads exactly n bytes from the connection with timeout.
func (c *GameConnection) Read(n int) ([]byte, error) {
	// Check if connection is shutting down
	select {
	case <-c.ctx.Done():
		return nil, ErrShuttingDown
	default:
	}

	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}

	if c.readTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return nil, fmt.Errorf("failed to set read deadline: %w", err)
		}
	}

	buf := make([]byte, n)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		// if timeout occurred, close connection to force reconnect
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			c.closeConn()
		}
		return nil, err
	}

	return buf, nil
}

// Write writes data to the connection with timeout.
func (c *GameConnection) Write(data []byte) (int, error) {
	select {
	case <-c.ctx.Done():
		return 0, ErrShuttingDown
	default:
	}

	conn, err := c.getConn()
	if err != nil {
		return 0, err
	}

	if c.writeTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return 0, fmt.Errorf("failed to set write deadline: %w", err)
		}
	}

	n, err := conn.Write(data)
	if err != nil {
		// if timeout occurred, close connection to force reconnect
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			c.closeConn()
		}
	}
	return n, err
}

// SetDeadline sets read and write deadlines.
func (c *GameConnection) SetDeadline(t time.Time) error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	return conn.SetDeadline(t)
}

// SetReadDeadline sets read deadline.
func (c *GameConnection) SetReadDeadline(t time.Time) error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	return conn.SetReadDeadline(t)
}

// SetWriteDeadline sets write deadline.
func (c *GameConnection) SetWriteDeadline(t time.Time) error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	return conn.SetWriteDeadline(t)
}

// Close gracefully shuts down the connection.
func (c *GameConnection) Close() error {
	c.cancel() // signal shutdown to all operations
	return c.closeConn()
}

// Connect establishes connection to endpoint with retry and backoff.
// Returns nil on success, otherwise an error.
// ctx allows to stop trying to connect.
func (c *GameConnection) Connect(ctx context.Context) error {
	connectCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		select {
		case <-connectCtx.Done():
			return context.Canceled
		case <-c.ctx.Done():
			return ErrShuttingDown
		default:
		}

		c.closeConn() // ignore possible error because we reestablish the connection anyway

		var conn net.Conn
		var err error

		conn, err = c.dialer.Dial(connectCtx, c.endpoint)

		if err == nil {
			// configure TCP connection
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				err = tcpConn.SetNoDelay(true)
				if err != nil {
					return err
				}
				err = tcpConn.SetKeepAlive(true)
				if err != nil {
					return err
				}
				err = tcpConn.SetKeepAlivePeriod(defaultKeepAlivePeriod)
				if err != nil {
					return err
				}
			}

			c.mu.Lock()
			c.conn = conn
			c.mu.Unlock()

			return nil
		}

		lastErr = err

		if attempt < c.maxRetries-1 {
			backoff := c.retryDelay * time.Duration(1<<uint(attempt))

			select {
			case <-connectCtx.Done():
				return context.Canceled
			case <-c.ctx.Done():
				return ErrShuttingDown
			case <-time.After(backoff):
				// continue to next attempt
			}
		}
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", c.maxRetries, lastErr)
}

// getConn safely returns the current connection or error.
func (c *GameConnection) getConn() (net.Conn, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return nil, ErrNotConnected
	}

	// check if connection is still alive
	var zero time.Time
	if err := conn.SetDeadline(zero); err != nil {
		c.closeConn()
		return nil, ErrNotConnected
	}

	return conn, nil
}

// closeConn safely closes and removes the current connection.
func (c *GameConnection) closeConn() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}

	return nil
}

// SetRetryDelay sets socket retry delay.
func (c *GameConnection) SetRetryDelay(delay time.Duration) {
	if delay > 0 {
		c.retryDelay = delay
	}
}

// SetMaxRetries sets socket max retries.
func (c *GameConnection) SetMaxRetries(maxRetries int) {
	if maxRetries > 0 {
		c.maxRetries = maxRetries
	}
}

// SetTimeouts sets read and write timeouts.
func (c *GameConnection) SetTimeouts(read, write time.Duration) {
	c.readTimeout = read
	c.writeTimeout = write
}
