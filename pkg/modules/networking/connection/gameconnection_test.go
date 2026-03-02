package connection

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/the-new-day/probogo/pkg/modules/networking/dial"
	"github.com/the-new-day/probogo/pkg/modules/networking/utils"
)

// --- Mocks ---

// mockConn implements net.Conn for testing.
type mockConn struct {
	readData    []byte
	readIndex   int
	writeData   []byte
	writeLock   sync.Mutex
	deadline    time.Time
	closeCalled bool
	readErr     error
	writeErr    error
	deadlineErr error
	onRead      func([]byte) (int, error) // For simulating partial reads
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	if m.onRead != nil {
		return m.onRead(b)
	}
	if m.readIndex >= len(m.readData) {
		return 0, io.EOF
	}
	n = copy(b, m.readData[m.readIndex:])
	m.readIndex += n
	return n, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.writeData = append(m.writeData, b...)
	return len(b), nil
}

func (m *mockConn) Close() error {
	m.closeCalled = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { m.deadline = t; return m.deadlineErr }
func (m *mockConn) SetReadDeadline(t time.Time) error  { m.deadline = t; return m.deadlineErr }
func (m *mockConn) SetWriteDeadline(t time.Time) error { m.deadline = t; return m.deadlineErr }

// mockDialer implements dial.Dialer for testing.
type mockDialer struct {
	conn            net.Conn
	err             error
	dialContextFunc func(ctx context.Context, addr string) (net.Conn, error)
}

func (m *mockDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	if m.dialContextFunc != nil {
		return m.dialContextFunc(ctx, addr)
	}
	return m.conn, m.err
}

// --- Helper Functions ---

func newTestGameConnection(dialer dial.Dialer) *GameConnection {
	return NewGameConnection(&utils.Address{Host: "test.com", Port: 1234}, dialer)
}

// --- Tests ---

func TestConnect_Success(t *testing.T) {
	mockConn := &mockConn{}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)

	err := gc.Connect(context.Background())
	require.NoError(t, err)
	assert.True(t, gc.IsConnected())
}

func TestConnect_FailureAndRetry(t *testing.T) {
	failCount := 0
	mockDialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, addr string) (net.Conn, error) {
			failCount++
			if failCount < 2 {
				return nil, errors.New("dial failed")
			}
			return &mockConn{}, nil
		},
	}
	gc := newTestGameConnection(mockDialer)
	gc.SetRetryDelay(1 * time.Millisecond) // Fast retries for test

	err := gc.Connect(context.Background())
	require.NoError(t, err)
	assert.True(t, gc.IsConnected())
	assert.Equal(t, 2, failCount)
}

func TestConnect_ContextCancel(t *testing.T) {
	mockDialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, addr string) (net.Conn, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}
	gc := newTestGameConnection(mockDialer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := gc.Connect(ctx)
	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, gc.IsConnected())
}

func TestConnect_MaxRetriesExceeded(t *testing.T) {
	expectedErr := errors.New("persistent dial error")
	mockDialer := &mockDialer{err: expectedErr}
	gc := newTestGameConnection(mockDialer)
	gc.SetRetryDelay(1 * time.Millisecond)
	gc.SetMaxRetries(3)

	err := gc.Connect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect after 3 attempts")
	assert.False(t, gc.IsConnected())
}

func TestRead_Success(t *testing.T) {
	mockConn := &mockConn{readData: []byte{1, 2, 3, 4}}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)
	require.NoError(t, gc.Connect(context.Background()))

	data, err := gc.Read(4)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3, 4}, data)
}

func TestRead_Timeout(t *testing.T) {
	mockConn := &mockConn{
		onRead: func(b []byte) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 0, netErrorWithTimeout{}
		},
	}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)
	require.NoError(t, gc.Connect(context.Background()))
	gc.SetTimeouts(1*time.Millisecond, 0) // Very short read timeout

	_, err := gc.Read(4)
	require.Error(t, err)
	// Check if the connection was closed due to timeout
	assert.False(t, gc.IsConnected())
}

func TestRead_NotConnected(t *testing.T) {
	gc := newTestGameConnection(&mockDialer{})
	_, err := gc.Read(4)
	assert.ErrorIs(t, err, ErrNotConnected)
}

func TestWrite_Success(t *testing.T) {
	mockConn := &mockConn{}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)
	require.NoError(t, gc.Connect(context.Background()))

	data := []byte{1, 2, 3}
	n, err := gc.Write(data)
	require.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, data, mockConn.writeData)
}

func TestWrite_Timeout(t *testing.T) {
	mockConn := &mockConn{
		writeErr: netErrorWithTimeout{},
	}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)
	require.NoError(t, gc.Connect(context.Background()))
	gc.SetTimeouts(0, 1*time.Millisecond) // Very short write timeout

	_, err := gc.Write([]byte{1, 2, 3})
	require.Error(t, err)
	assert.False(t, gc.IsConnected())
}

func TestClose(t *testing.T) {
	mockConn := &mockConn{}
	mockDialer := &mockDialer{conn: mockConn}
	gc := newTestGameConnection(mockDialer)
	require.NoError(t, gc.Connect(context.Background()))

	err := gc.Close()
	require.NoError(t, err)
	assert.True(t, mockConn.closeCalled)
	assert.False(t, gc.IsConnected())

	// Subsequent operations should return ErrShuttingDown or ErrNotConnected
	_, err = gc.Read(1)
	assert.ErrorIs(t, err, ErrShuttingDown)
}

// Helper for timeout error simulation
type netErrorWithTimeout struct{}

func (e netErrorWithTimeout) Error() string   { return "timeout" }
func (e netErrorWithTimeout) Timeout() bool   { return true }
func (e netErrorWithTimeout) Temporary() bool { return false }
