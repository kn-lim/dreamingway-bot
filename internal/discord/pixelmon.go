package discord

import "github.com/bwmarrin/discordgo"

func status(*discordgo.Interaction) (discordgo.InteractionResponse, error) {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Status!",
		},
	}, nil
}

func test(*discordgo.Interaction) (discordgo.InteractionResponse, error) {
	return discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Test!",
		},
	}, nil
}
