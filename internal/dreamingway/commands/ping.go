package commands

import "github.com/bwmarrin/discordgo"

// ping returns "Pong!"
func ping(i *discordgo.Interaction) (string, error) {
	return "Pong!", nil
}
