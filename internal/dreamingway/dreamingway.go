package dreamingway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

const (
	// Discord Webhook Base URL
	WEBHOOK_BASE_URL = "https://discord.com/api"
)

type Dreamingway interface {
	DeferredMessage() discordgo.InteractionResponse
	SendDeferredMessage(appID, token, content string) error
}

type DreamingwayBot struct {
	Client *discordgo.Session
}

// NewDreamingway creates a new DreamingwayBot instance
func NewDreamingway(token string) (*DreamingwayBot, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Errorw("failed to create discord client",
				"error", err,
			)
		}
		return nil, err
	}

	return &DreamingwayBot{
		Client: client,
	}, nil
}

// DeferredMessage returns a deferred message response
func (d *DreamingwayBot) DeferredMessage() discordgo.InteractionResponse {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}
}

// SendDeferredMessage sends a deferred message to a Discord channel
func (d *DreamingwayBot) SendDeferredMessage(appID, token, content string) error {
	payload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Errorw("failed to marshal payload",
				"error", err,
			)
		}
		return err
	}

	url := fmt.Sprintf("%v/v%v/webhooks/%v/%v", WEBHOOK_BASE_URL, os.Getenv("DISCORD_API_VERSION"), appID, token)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Errorw("failed to create http POST request",
				"error", err,
			)
		}
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Errorw("failed to send http POST request",
				"error", err,
			)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			if utils.Logger != nil {
				utils.Logger.Errorw("failed to decode response body",
					"error", err,
				)
			}
		}
		if utils.Logger != nil {
			utils.Logger.Errorw("discord API error",
				"error", result,
			)
		}
		return fmt.Errorf("discord API error: %v", result)
	}

	return nil
}

// GetServerName returns the server name from a guildID
func (d *DreamingwayBot) GetServerName(guildID string) (string, error) {
	// Get server name from guild_id
	guild, err := d.Client.Guild(guildID)
	if err != nil {
		utils.Logger.Errorw("error getting guild",
			"guild_id", guildID,
			"error", err,
		)
		return "", err
	}

	return guild.Name, nil
}
