package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/the-new-day/probogo/pkg/modules/networking/connection"
	"github.com/the-new-day/probogo/pkg/modules/networking/dial"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/modules/proxy"
	"github.com/the-new-day/probogo/pkg/packets"

	_ "github.com/the-new-day/probogo/pkg/packets/battleinfo"
	_ "github.com/the-new-day/probogo/pkg/packets/battlemechanics"
	_ "github.com/the-new-day/probogo/pkg/packets/chat"
	_ "github.com/the-new-day/probogo/pkg/packets/entry"
	_ "github.com/the-new-day/probogo/pkg/packets/garage"
	_ "github.com/the-new-day/probogo/pkg/packets/lobby"
	_ "github.com/the-new-day/probogo/pkg/packets/network"
	_ "github.com/the-new-day/probogo/pkg/packets/others"
	_ "github.com/the-new-day/probogo/pkg/packets/shop"
	_ "github.com/the-new-day/probogo/pkg/packets/status"
	_ "github.com/the-new-day/probogo/pkg/packets/turrets"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytesSent  int64
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesSent += int64(n)
	return n, err
}

type LauncherProxy struct {
	httpListenAddr string
	tcpListenAddr  string
	serverAddr     string
	packetRegistry *packets.PacketRegistry

	mu         sync.Mutex
	httpServer *http.Server
}

func NewLauncherProxy(httpAddr, tcpAddr, serverAddr string, registry *packets.PacketRegistry) *LauncherProxy {
	return &LauncherProxy{
		httpListenAddr: httpAddr,
		tcpListenAddr:  tcpAddr,
		serverAddr:     serverAddr,
		packetRegistry: registry,
	}
}

func (lp *LauncherProxy) Run(ctx context.Context) error {
	log.Printf("Starting proxy: HTTP=%s TCP=%s → Server=%s", lp.httpListenAddr, lp.tcpListenAddr, lp.serverAddr)
	go lp.runGameProxy(ctx)
	return lp.runHTTPServer(ctx)
}

func (lp *LauncherProxy) runGameProxy(ctx context.Context) {
	listener, err := net.Listen("tcp", lp.tcpListenAddr)
	if err != nil {
		log.Printf("TCP proxy failed: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("TCP proxy listening on %s", lp.tcpListenAddr)

	var connCount int
	for {
		clientRaw, err := listener.Accept()
		if err != nil {
			continue
		}

		connCount++
		go lp.handleClient(clientRaw, connCount)
	}
}

func (lp *LauncherProxy) handleClient(clientRaw net.Conn, connID int) {
	defer clientRaw.Close()

	clientConn := connection.NewEstablishedConnection(clientRaw)

	serverConn := connection.NewGameConnection(
		lp.serverAddr,
		dial.NewDirectDialer(10*time.Second),
	)

	if err := serverConn.Connect(context.Background()); err != nil {
		log.Printf("Client %d: connection failed: %v", connID, err)
		return
	}
	defer serverConn.Close()

	serverProt := protection.NewXorProtection(false)
	clientProt := protection.NewXorProtection(true)

	gameProxy := proxy.NewProxy(
		serverConn,
		serverProt,
		clientConn,
		clientProt,
		lp.packetRegistry,
	)

	log.Printf("Client %d connected", connID)
	gameProxy.Run(context.Background())
	log.Printf("Client %d disconnected", connID)
}

func (lp *LauncherProxy) runHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./prelauncher"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		fs.ServeHTTP(lrw, r)
		log.Printf("HTTP %s %s from %s - %d bytes", r.Method, r.URL.Path, r.RemoteAddr, lrw.bytesSent)
	}))

	lp.httpServer = &http.Server{Addr: lp.httpListenAddr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		if err := lp.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return lp.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func (lp *LauncherProxy) Close() error {
	if lp.httpServer != nil {
		return lp.httpServer.Close()
	}
	return nil
}

func main() {
	httpAddr := flag.String("http", "127.0.0.1:1211", "HTTP server address")
	tcpAddr := flag.String("tcp", "127.0.0.1:1213", "TCP proxy address")
	serverAddr := flag.String("server", "146.59.110.146:25565", "Game server address")
	flag.Parse()

	registry := packets.GlobalPacketRegistry()

	proxy := NewLauncherProxy(*httpAddr, *tcpAddr, *serverAddr, registry)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Printf("Launching proxy - HTTP:%s TCP:%s → %s", *httpAddr, *tcpAddr, *serverAddr)

	if err := proxy.Run(ctx); err != nil {
		log.Fatalf("Proxy error: %v", err)
	}
}
