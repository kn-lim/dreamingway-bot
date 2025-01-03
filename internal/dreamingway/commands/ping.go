package commands

import "github.com/bwmarrin/discordgo"

func ping(i *discordgo.Interaction) (string, error) {
	return "Pong!", nil
}
