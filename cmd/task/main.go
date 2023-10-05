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
	log.Printf("Received interaction: %+v\n", interaction)

	application_id := interaction.AppID
	interaction_token := interaction.Token

	url := fmt.Sprintf("%v/v%v/webhooks/%v/%v", discord.DiscordBaseURL, os.Getenv("DISCORD_API_VERSION"), application_id, interaction_token)

	log.Printf("Discord API URL: %s", url)

	payloadBytes, err := json.Marshal(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This is from the task lambda!",
		},
	})
	if err != nil {
		log.Fatalf("Error! Couldn't marshal JSON: %v", err)
	}

	log.Printf("Sending payload: %s", string(payloadBytes))

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
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
			log.Printf("Error! Couldn't decode result: %v", err)
			return err
		}
		log.Printf("Error! Discord API Error: %v", result)
		return fmt.Errorf("discord API error: %v", result)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
