package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

const DiscordBaseURL = "https://discord.com/api"

func DeferredMessage(i *discordgo.Interaction) (discordgo.InteractionResponse, error) {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}, nil
}

func SendDeferredMessage(appID string, token string, content string) error {
	log.Printf("Sending message: %s", content)

	payload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %v", err)
	}

	url := fmt.Sprintf("%v/v%v/webhooks/%v/%v", DiscordBaseURL, os.Getenv("DISCORD_API_VERSION"), appID, token)
	// log.Printf("Discord API URL: %s", url)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("couldn't create http request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bot "+os.Getenv("DISCORD_BOT_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("couldn't send the http request: %v", err)
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
