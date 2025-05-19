package dreamingway

import (
	"github.com/disgoorg/disgo/discord"
)

// func GetUsername(interaction discordgo.Interaction) string {
// 	username := "???"
// 	if interaction.Member != nil {
// 		username = interaction.Member.User.Username
// 	} else if interaction.User != nil {
// 		username = interaction.User.Username
// 	}
// 	return username
// }

func GetDeferredMessageResponse() discord.InteractionResponse {
	return discord.InteractionResponse{
		Type: discord.InteractionResponseTypeDeferredCreateMessage,
	}
}
