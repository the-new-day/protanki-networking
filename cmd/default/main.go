package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/the-new-day/protanki-networking/pkg/modules/proxy"
	"github.com/the-new-day/protanki-networking/pkg/packets"
	"github.com/the-new-day/protanki-networking/pkg/packets/chat"
)

func main() {
	httpAddr := flag.String("http", "127.0.0.1:1211", "HTTP listen address")
	tcpAddr := flag.String("tcp", "127.0.0.1:1213", "TCP proxy address")
	serverAddr := flag.String("server", "146.59.110.146:25565", "Game server address")
	flag.Parse()

	registry := packets.GlobalPacketRegistry()

	launcher := NewProxyLauncher(*httpAddr, *tcpAddr, *serverAddr, registry)

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
		case proxy, ok := <-launcher.ProxyCh:
			if !ok {
				return
			}

			setupProxyHandlers(proxy)
		}
	}
}

func setupProxyHandlers(proxy *proxy.Proxy) {
	proxy.OnServerToClient(func(packet packets.Packet) packets.Packet {
		if packet.ID() == packets.ReceiveLobbyChatID {
			messages := packets.Attr[[]map[string]any]("messages", packet)

			if len(messages) != 1 {
				return packet
			}

			for _, message := range messages {
				if strings.Contains(strings.ToLower(message["text"].(string)), "волки вотана") {
					fake := chat.NewReceiveLobbyChatPacket()
					attributes := map[string]any{
						"authorStatus":  packets.Boolshortern(),
						"systemMessage": true,
						"targetStatus":  packets.Boolshortern(),
						"text":          "Волки Вотана MENTIONED",
						"warning":       false,
					}

					fake.UnwrapValues([]map[string]any{attributes})

					err := proxy.SendToClient(packet)
					if err != nil {
						log.Printf("[ERROR]: %s", err)
					}

					err = proxy.SendToClient(fake)
					if err != nil {
						log.Printf("[ERROR]: %s", err)
					}

					return nil
				}
			}
		}
		return packet
	})
}
