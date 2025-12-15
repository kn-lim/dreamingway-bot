package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/projectzomboid"
	"github.com/kn-lim/chattingway/rcon"
	"github.com/kn-lim/dreamingway-bot/internal/constants"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
)

const (
	PLAYER_COUNT_REGEX = `Players connected \((\d+)\):`
)

// Project Zomboid related commands
func pz(i discord.Interaction) (string, error) {
	// Validate user has permission to run admin commands
	member := i.Member()
	if member == nil {
		return "", errors.New("command must be sent from a server")
	}
	roles, err := dreamingway.GetAllGuildRoles(i.GuildID(), os.Getenv("DISCORD_API_VERSION"), os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return "", err
	}

	is_admin := false
	for _, role := range roles {
		if role.Name == os.Getenv("PZ_DISCORD_ADMIN_ROLE") {
			is_admin = slices.Contains(member.RoleIDs, role.ID)
			break
		}
	}

	// Handle subcommand
	subCommand := i.(discord.ApplicationCommandInteraction).SlashCommandInteractionData()
	switch *subCommand.SubCommandName {
	case "rcon": // /pz rcon <command>
		if !is_admin {
			return constants.UNAUTHORIZED, nil
		}

		output, err := rcon.Run(os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PORT"), os.Getenv("PZ_RCON_PASSWORD"), subCommand.String("command"))
		if err != nil {
			return "", err
		}

		if output == "" {
			return "Successfully sent the RCON command", nil
		} else {
			return fmt.Sprintf("Successfully sent the RCON command and received the output: `%s`", output), nil
		}
	case "start": // /pz start
		if !is_admin {
			return constants.UNAUTHORIZED, nil
		}

		if err := projectzomboid.Start(context.TODO(), os.Getenv("PZ_HOST_INSTANCE_ID"), os.Getenv("PZ_HOST_REGION"), os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PORT"), os.Getenv("PZ_RCON_PASSWORD")); err != nil {
			return "", err
		}

		return constants.PZ_STATUS_ONLINE, nil
	case "status": // /pz status
		status, err := projectzomboid.Status(os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PORT"), os.Getenv("PZ_RCON_PASSWORD"))
		if err != nil {
			return "", err
		}

		if status {
			return constants.PZ_STATUS_ONLINE, nil
		} else {
			return constants.PZ_STATUS_OFFLINE, nil
		}
	case "stop": // /pz stop
		if !is_admin {
			return constants.UNAUTHORIZED, nil
		}

		if err := projectzomboid.Stop(context.TODO(), os.Getenv("PZ_HOST_INSTANCE_ID"), os.Getenv("PZ_HOST_REGION"), os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PORT"), os.Getenv("PZ_RCON_PASSWORD")); err != nil {
			return "", err
		}

		return constants.PZ_STATUS_OFFLINE, nil
	}

	return "", errors.New("invalid pz subcommand")
}
