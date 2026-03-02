package dial

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

// UpgradeToTLSIfNeeded upgrades connection to TLS if the target port is 443.
func UpgradeToTLSIfNeeded(ctx context.Context, conn net.Conn, targetAddr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(targetAddr)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("invalid target address: %w", err)
	}

	if port != "443" {
		return conn, nil
	}

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: host,
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- tlsConn.Handshake()
	}()

	select {
	case <-ctx.Done():
		conn.Close()
		return nil, ctx.Err()
	case err := <-errCh:
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("TLS handshake failed: %w", err)
		}
	}

	return tlsConn, nil
}
