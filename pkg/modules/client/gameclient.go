package client

import (
	"context"

	"github.com/the-new-day/probogo/pkg/modules/networking/utils"
	"github.com/the-new-day/probogo/pkg/modules/protection"
)

// GameClient manages connection with the server and serves as
// a high-level facade for sending/receiving packets and handling them.
type GameClient struct {
	name            string
	gameCredentials utils.Credentials
	proxy           *utils.ProxyCredentials

	emergencyHalt context.Context
	cancel        context.CancelFunc

	protection protection.Protection
}

func NewGameClient(
	name string,
	serverAddress *utils.Address,
	proxy *utils.ProxyCredentials,
	gameCredentials utils.Credentials,
	protection protection.Protection,
) *GameClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &GameClient{
		name:            name,
		gameCredentials: gameCredentials,
		proxy:           nil,

		emergencyHalt: ctx,
		cancel:        cancel,

		protection: protection,
	}
}
