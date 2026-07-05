package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/the-new-day/protanki-networking/pkg/modules/networking"
	"github.com/the-new-day/protanki-networking/pkg/modules/proxy"
	"github.com/the-new-day/protanki-networking/pkg/packets"
	"github.com/the-new-day/protanki-networking/pkg/utils"
)

func main() {
	httpAddr := flag.String("http", "127.0.0.1:1211", "HTTP listen address")
	tcpAddr := flag.String("tcp", "127.0.0.1:1213", "TCP proxy address")
	serverAddr := flag.String("server", "146.59.110.146:25565", "Game server address")
	flag.Parse()

	registry := packets.GlobalPacketRegistry()

	launcher := proxy.NewProxyLauncher(
		*httpAddr,
		*tcpAddr,
		*serverAddr,
		registry,
		proxy.WithLogger(log.Default()),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := launcher.Run(ctx); err != nil {
			log.Fatalf("Proxy stopped with error: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		case proxy, ok := <-launcher.ProxyCh():
			if !ok {
				return
			}

			setupProxyHandlers(proxy)
		}
	}
}

func setupProxyHandlers(proxy *proxy.Proxy) {
	proxy.OnServerToClient(func(packet packets.Packet) packets.Packet {
		log.Printf("[SERVER]: %s\nPayload:\n%s", packets.GetName(packet.ID()), utils.ShortView(packet.Data(), 20))
		return packet
	})
	proxy.OnClientToServer(func(packet packets.Packet) packets.Packet {
		log.Printf("[CLIENT]: %s\nPayload:\n%s", packets.GetName(packet.ID()), utils.ShortView(packet.Data(), 20))
		return packet
	})
	proxy.OnServerError(func(pr networking.PacketResult) {
		log.Printf("[ERROR(S->C)]: Packet=%s | Err: %s\n", packets.GetName(pr.ID), pr.Err)
	})
	proxy.OnClientError(func(pr networking.PacketResult) {
		log.Printf("[ERROR(C->S)]: Packet=%s | Err: %s\n", packets.GetName(pr.ID), pr.Err)
	})
}
