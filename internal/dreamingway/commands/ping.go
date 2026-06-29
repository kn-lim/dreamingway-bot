package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/v2/healthcheck"
)

// ping returns "Pong!"
func ping(i discord.Interaction) (string, error) {
	return healthcheck.Ping(), nil
}
