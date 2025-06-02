package commands

import "github.com/disgoorg/disgo/discord"

// ping returns "Pong!"
func ping(i discord.Interaction) (string, error) {
	return "Pong!", nil
}
