package discord

import "github.com/bwmarrin/discordgo"

const DiscordBaseURL = "https://discord.com/api"

func DeferredMessage(i *discordgo.Interaction) (discordgo.InteractionResponse, error) {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}, nil
}
