package dial

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, targetAddr string) (net.Conn, error)
}
