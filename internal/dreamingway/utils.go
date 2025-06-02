package dreamingway

import (
	"github.com/disgoorg/disgo/discord"
)

func GetDeferredMessageResponse() discord.InteractionResponse {
	return discord.InteractionResponse{
		Type: discord.InteractionResponseTypeDeferredCreateMessage,
	}
}
