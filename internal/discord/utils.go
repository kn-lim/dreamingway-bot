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

const (
	DiscordBaseURL = "https://discord.com/api"

	ErrMissingRole = "You don't have the required role to use this command!"
)

type options struct {
	client *http.Client
	url    string
}
type Option func(*options)

func WithClient(client *http.Client) Option {
	return func(o *options) {
		o.client = client
	}
}

func WithURL(url string) Option {
	return func(o *options) {
		o.url = url
	}
}

func CheckRole(roles []string, requiredRole string) bool {
	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}

	return false
}

func DeferredMessage() (discordgo.InteractionResponse, error) {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}, nil
}

func SendDeferredMessage(appID string, token string, content string, opts ...Option) error {
	log.Printf("Sending message: %s", content)

	// Defaults
	config := &options{
		client: &http.Client{},
		url:    DiscordBaseURL,
	}
	for _, opt := range opts {
		opt(config)
	}

	payload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %v", err)
	}

	url := fmt.Sprintf("%v/v%v/webhooks/%v/%v", config.url, os.Getenv("DISCORD_API_VERSION"), appID, token)
	// log.Printf("Discord API URL: %s", url)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error! couldn't create http request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bot "+os.Getenv("DISCORD_BOT_TOKEN"))

	client := config.client
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error! couldn't send the http request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			log.Printf("Error! Couldn't decode result: %v", err)
			return err
		}
		log.Printf("Error! Discord API Error: %v", result)
		return fmt.Errorf("error! discord API error: %v", result)
	}

	return nil
}
