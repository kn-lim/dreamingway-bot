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

type message struct {
	Content    string     `json:"content"`
	Components components `json:"components"`
}

type components struct {
	Type       int      `json:"type"`
	Components []button `json:"components"`
}

type button struct {
	Type     int    `json:"type"`
	Label    string `json:"label"`
	Style    int    `json:"style"`
	CustomID string `json:"custom_id"`
}

type options struct {
	actionsRow *discordgo.ActionsRow

	// For tests
	client *http.Client
	url    string
}
type Option func(*options)

func WithActionsRow(actionsRow *discordgo.ActionsRow) Option {
	return func(o *options) {
		o.actionsRow = actionsRow
	}
}

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
		// For tests
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

func SendDeferredMessageWithComponents(appID string, token string, content string, opts ...Option) error {
	log.Printf("Sending message: %s", content)

	// Defaults
	config := &options{
		// For tests
		client: &http.Client{},
		url:    DiscordBaseURL,
	}
	for _, opt := range opts {
		opt(config)
	}

	message := message{
		Content:    content,
		Components: components{},
	}
	if config.actionsRow != nil {
		var buttons []button
		for _, i := range config.actionsRow.Components {
			buttons = append(buttons, button{
				Type:     int(i.Type()),
				Label:    i.(discordgo.Button).Label,
				Style:    int(i.(discordgo.Button).Style),
				CustomID: i.(discordgo.Button).CustomID,
			})
		}

		message.Components = components{
			Type:       int(discordgo.ActionsRowComponent),
			Components: buttons,
		}
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %v", err)
	}

	log.Printf("Payload: %v", payload)

	temp, _ := json.Marshal(map[string]string{
		"content": "test",
	})
	log.Printf("Correct Payload: %v", temp)

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
