package proxy

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/the-new-day/protanki-networking/pkg/modules/networking/connection"
	"github.com/the-new-day/protanki-networking/pkg/modules/networking/dial"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Logger defines the minimal logging interface used by ProxyLauncher.
type Logger interface {
	Printf(format string, v ...any)
}

// ProxyLauncherOption configures a ProxyLauncher.
type ProxyLauncherOption func(*ProxyLauncher)

// WithLogger sets a custom logger. If nil, logging is disabled.
func WithLogger(logger Logger) ProxyLauncherOption {
	return func(pl *ProxyLauncher) {
		pl.logger = logger
	}
}

// WithDialer sets a custom dialer for connecting to the upstream server.
func WithDialer(d dial.Dialer) ProxyLauncherOption {
	return func(pl *ProxyLauncher) {
		pl.dialer = d
	}
}

// WithServerProtectionFactory sets a factory function that creates a new
// Protection instance for the server side of each proxy session.
func WithServerProtectionFactory(factory func() protection.Protection) ProxyLauncherOption {
	return func(pl *ProxyLauncher) {
		pl.serverProtectionFactory = factory
	}
}

// WithClientProtectionFactory sets a factory function that creates a new
// Protection instance for the client side of each proxy session.
func WithClientProtectionFactory(factory func() protection.Protection) ProxyLauncherOption {
	return func(pl *ProxyLauncher) {
		pl.clientProtectionFactory = factory
	}
}

// ProxyLauncher manages both HTTP and TCP proxy servers.
type ProxyLauncher struct {
	httpAddr   string
	tcpAddr    string
	serverAddr string
	registry   *packets.PacketRegistry

	// ProxyCh is a channel that delivers new Proxy instances as they are created.
	// Each new client connection to the TCP proxy will result in a new Proxy being
	// sent on this channel. The channel is closed when the ProxyLauncher's Run
	// method exits.
	//
	// This allows consumers to:
	//   - Access and control each proxy session individually
	//   - Monitor new connections as they happen
	//   - Perform additional setup or logging per connection
	//   - Implement connection limiting or filtering
	//
	// Example usage:
	//
	//	go func() {
	//	    for proxy := range launcher.ProxyCh {
	//	        log.Printf("New proxy session started")
	//	        // Store proxy reference, monitor it, etc.
	//	    }
	//	}()
	ProxyCh chan *Proxy

	logger                  Logger
	dialer                  dial.Dialer
	serverProtectionFactory func() protection.Protection
	clientProtectionFactory func() protection.Protection
}

// NewProxyLauncher creates a new ProxyLauncher with the given addresses and registry.
// Additional options can be provided to customize logging, dialer, and protection factories.
// By default:
//   - logging is disabled
//   - dialer is dial.NewDirectDialer(10*time.Second)
//   - server protection is protection.NewXorProtection(false)
//   - client protection is protection.NewXorProtection(true)
func NewProxyLauncher(httpAddr, tcpAddr, serverAddr string, registry *packets.PacketRegistry, opts ...ProxyLauncherOption) *ProxyLauncher {
	pl := &ProxyLauncher{
		httpAddr:   httpAddr,
		tcpAddr:    tcpAddr,
		serverAddr: serverAddr,
		registry:   registry,
		ProxyCh:    make(chan *Proxy),

		// defaults
		logger: nil,
		dialer: dial.NewDirectDialer(10 * time.Second),
		serverProtectionFactory: func() protection.Protection {
			return protection.NewXorProtection(false)
		},
		clientProtectionFactory: func() protection.Protection {
			return protection.NewXorProtection(true)
		},
	}

	for _, opt := range opts {
		opt(pl)
	}

	return pl
}

// logf prints a formatted message if a logger is set.
func (p *ProxyLauncher) logf(format string, v ...any) {
	if p.logger != nil {
		p.logger.Printf(format, v...)
	}
}

// Run starts both HTTP and TCP servers. It blocks until ctx is cancelled.
func (p *ProxyLauncher) Run(ctx context.Context) error {
	p.logf("Starting: HTTP=%s TCP=%s → %s", p.httpAddr, p.tcpAddr, p.serverAddr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := p.runHTTP(ctx); err != nil {
			p.logf("HTTP error: %v", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := p.runTCP(ctx); err != nil {
			p.logf("TCP error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	wg.Wait()
	close(p.ProxyCh)
	return nil
}

func (p *ProxyLauncher) runTCP(ctx context.Context) error {
	listener, err := connection.NewClientListener(p.tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	var id int
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			client, err := listener.Accept()
			if err != nil {
				continue
			}
			id++
			go p.handleClient(ctx, client, id)
		}
	}
}

func (p *ProxyLauncher) handleClient(ctx context.Context, client *connection.ClientConnection, id int) {
	defer client.Close()

	server := connection.NewServerConnection(
		p.serverAddr,
		p.dialer,
	)

	if err := server.Connect(ctx); err != nil {
		p.logf("Client %d: connect failed: %v", id, err)
		return
	}
	defer server.Close()

	gameProxy := NewProxy(
		server,
		p.serverProtectionFactory(),
		client,
		p.clientProtectionFactory(),
		p.registry,
	)

	p.logf("Client %d connected", id)

	select {
	case p.ProxyCh <- gameProxy:
	default:
	}

	gameProxy.Run(ctx)
}

func (p *ProxyLauncher) runHTTP(ctx context.Context) error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./prelauncher"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
		p.logf("HTTP %s %s", r.Method, r.URL.Path)
	}))

	server := &http.Server{Addr: p.httpAddr, Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	p.logf("HTTP listening on %s", p.httpAddr)

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
