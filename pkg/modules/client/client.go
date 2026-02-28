package client

import (
	"context"

	"github.com/the-new-day/probogo/pkg/modules/networking"
	"github.com/the-new-day/probogo/pkg/modules/processing"
	"github.com/the-new-day/probogo/pkg/modules/protection"
)

// Client manages connection with the server and serves as
// a high-level facade for sending/receiving packets and handling them.
type Client struct {
	name            string
	gameCredentials networking.Credentials
	proxy           *networking.ProxyCredentials

	emergencyHalt context.Context
	cancel        context.CancelFunc

	protection protection.Protection
	socket     *networking.GameSocket

	packetHandler processing.PacketHandler
}

func NewClient(
	name string,
	serverAddress *networking.Address,
	proxy *networking.ProxyCredentials,
	gameCredentials networking.Credentials,
	protection protection.Protection,
	packetHandler processing.PacketHandler,
) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	socket := networking.NewGameSocket(serverAddress, protection, proxy, ctx)

	return &Client{
		name:            name,
		gameCredentials: gameCredentials,
		proxy:           nil,

		emergencyHalt: ctx,
		cancel:        cancel,

		protection: protection,
		socket:     socket,

		packetHandler: packetHandler,
	}
}
