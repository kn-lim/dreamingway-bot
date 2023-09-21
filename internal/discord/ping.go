package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ping(*discordgo.Interaction) (discordgo.InteractionResponse, error) {
	log.Print("Pong!")

	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	}, nil
}
