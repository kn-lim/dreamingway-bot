package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

func handler(interaction discordgo.Interaction) error {
	application_id := interaction.AppID
	interaction_token := interaction.Token
	message_id := interaction.Message.ID

	url := fmt.Sprintf("%v/%v/webhooks/%v/%v/messages/%v", discord.DiscordBaseURL, os.Getenv("DISCORD_API_VERSION"), application_id, interaction_token, message_id)

	log.Println(url)

	payloadBytes, err := json.Marshal(discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This is from the task lambda!",
		},
	})
	if err != nil {
		log.Fatalf("Error! Couldn't marshal JSON: %v", err)
	}

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error! Couldn't create http request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bot "+os.Getenv("DISCORD_BOT_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error! Couldn't send the http request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			log.Fatalf("Error! Couldn't decode result: %v", err)
		}
		log.Fatalf("Error! Discord API Error: %v", result)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
