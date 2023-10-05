package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

func handler(interaction discordgo.Interaction) error {
	application_id := interaction.AppID
	interaction_token := interaction.Token

	url := fmt.Sprintf("%v/%v/webhooks/%v/%v/messages/@original", discord.DiscordBaseURL, os.Getenv("DISCORD_API_VERSION"), application_id, interaction_token)

	payloadBytes, _ := json.Marshal(discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This is from the task lambda!",
		},
	})

	resp, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadBytes))
	defer resp.Body.Close()

	return nil
}

func main() {
	lambda.Start(handler)
}
