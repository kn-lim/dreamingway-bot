package dreamingway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

// GetDeferredMessageResponse returns a InteractionResponseTypeDeferredCreateMessage interaction response
func GetDeferredMessageResponse() discord.InteractionResponse {
	return discord.InteractionResponse{
		Type: discord.InteractionResponseTypeDeferredCreateMessage,
	}
}

// GetAllGuildRoles returns all roles present in a given Discord guild
func GetAllGuildRoles(guildID *snowflake.ID, apiVersion, token string) ([]discord.Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v%s/guilds/%s/roles", WEBHOOK_BASE_URL, apiVersion, guildID.String()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Errorw("failed to send http GET request",
				"error", err,
			)
		}

		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
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
		return nil, fmt.Errorf("discord API error: %v", result)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var roles []discord.Role
	if err := json.Unmarshal(body, &roles); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return roles, nil
}
