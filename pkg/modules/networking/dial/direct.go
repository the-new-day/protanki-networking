package dial

import (
	"context"
	"fmt"
	"net"
	"time"
)

type DirectDialer struct {
	timeout time.Duration
}

func NewDirectStrategy(timeout time.Duration) *DirectDialer {
	return &DirectDialer{timeout: timeout}
}

func (s *DirectDialer) Dial(ctx context.Context, targetAddr string) (net.Conn, error) {
	var d net.Dialer
	if s.timeout > 0 {
		d.Timeout = s.timeout
	}

	conn, err := d.DialContext(ctx, "tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("direct dial failed: %w", err)
	}

	return UpgradeToTLSIfNeeded(ctx, conn, targetAddr)
}
