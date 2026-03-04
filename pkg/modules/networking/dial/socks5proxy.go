package dial

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

// Socks5ProxyDialer implements dialing through SOCKS5 proxy.
type Socks5ProxyDialer struct {
	proxyAddr string
	auth      *proxy.Auth
	timeout   time.Duration
}

func NewSocks5Dialer(proxyAddr string, username, password string, timeout time.Duration) *Socks5ProxyDialer {
	var auth *proxy.Auth
	if username != "" {
		auth = &proxy.Auth{
			User:     username,
			Password: password,
		}
	}

	return &Socks5ProxyDialer{
		proxyAddr: proxyAddr,
		auth:      auth,
		timeout:   timeout,
	}
}

func (s *Socks5ProxyDialer) Dial(ctx context.Context, targetAddr string) (net.Conn, error) {
	dialer := &socks5Dialer{
		proxyAddr: s.proxyAddr,
		auth:      s.auth,
		timeout:   s.timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("SOCKS5 dial failed: %w", err)
	}

	return UpgradeToTLSIfNeeded(ctx, conn, targetAddr)
}

// socks5Dialer wrappers for SOCKS5 proxy with context support.
type socks5Dialer struct {
	proxyAddr string
	auth      *proxy.Auth
	timeout   time.Duration
}

func (d *socks5Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5(network, d.proxyAddr, d.auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return nil, errors.New("SOCKS5 dialer does not support ContextDialer")
	}

	if d.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d.timeout)
		defer cancel()
	}

	return contextDialer.DialContext(ctx, network, addr)
}
