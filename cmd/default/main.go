package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

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
	wasFakeSent := false

	proxy.OnServerToClient(func(packet packets.Packet) packets.Packet {
		if packet.ID() == packets.ReceiveLobbyChatID {
			messages := packets.Attr[[]map[string]any]("messages", packet)

			for _, message := range messages {
				authorStatus := message["authorStatus"].(map[string]any)
				targetStatus := message["targetStatus"].(map[string]any)

				authorUsername := authorStatus["username"]
				targetUsername := targetStatus["username"]

				if authorUsername == "slowed" {
					authorStatus["modLevel"] = int32(4)
				} else if targetUsername == "slowed" {
					targetStatus["modLevel"] = int32(4)
				}
			}

			packet.Set("messages", messages)

			if !wasFakeSent {
				fake := chat.NewReceiveLobbyChatPacket()
				attributes := map[string]any{
					"authorStatus":  packets.Boolshortern(),
					"systemMessage": true,
					"targetStatus":  packets.Boolshortern(),
					"text":          "Hello World",
					"warning":       false,
				}
				fake.UnwrapValues([]map[string]any{attributes})

				err := proxy.SendToClient(fake)
				if err == nil {
					wasFakeSent = true
				}
			}
		}
		return packet
	})
}
