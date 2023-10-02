package discord

import (
	"github.com/bwmarrin/discordgo"
)

func ping(*discordgo.Interaction) (discordgo.InteractionResponse, error) {
	// log.Println("ping")

	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	}, nil
}
