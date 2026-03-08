package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/the-new-day/protanki-networking/pkg/modules/networking/connection"
	"github.com/the-new-day/protanki-networking/pkg/modules/networking/dial"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
	"github.com/the-new-day/protanki-networking/pkg/modules/proxy"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

type ProxyLauncher struct {
	httpAddr   string
	tcpAddr    string
	serverAddr string
	registry   *packets.PacketRegistry

	ProxyCh chan *proxy.Proxy
}

func NewProxyLauncher(httpAddr, tcpAddr, serverAddr string, registry *packets.PacketRegistry) *ProxyLauncher {
	return &ProxyLauncher{
		httpAddr:   httpAddr,
		tcpAddr:    tcpAddr,
		serverAddr: serverAddr,
		registry:   registry,
		ProxyCh:    make(chan *proxy.Proxy),
	}
}

func (p *ProxyLauncher) Run(ctx context.Context) error {
	log.Printf("Starting: HTTP=%s TCP=%s → %s", p.httpAddr, p.tcpAddr, p.serverAddr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := p.runHTTP(ctx); err != nil {
			log.Printf("HTTP error: %v", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := p.runTCP(ctx); err != nil {
			log.Printf("TCP error: %v", err)
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

	log.Printf("TCP listening on %s", p.tcpAddr)

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
		dial.NewDirectDialer(10*time.Second),
	)

	if err := server.Connect(ctx); err != nil {
		log.Printf("Client %d: connect failed: %v", id, err)
		return
	}
	defer server.Close()

	gameProxy := proxy.NewProxy(
		server,
		protection.NewXorProtection(false),
		client,
		protection.NewXorProtection(true),
		p.registry,
	)

	log.Printf("Client %d connected", id)

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
		log.Printf("HTTP %s %s", r.Method, r.URL.Path)
	}))

	server := &http.Server{Addr: p.httpAddr, Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	log.Printf("HTTP listening on %s", p.httpAddr)

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
