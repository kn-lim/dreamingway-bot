package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/projectzomboid"
	"github.com/kn-lim/chattingway/rcon"
	"github.com/kn-lim/dreamingway-bot/internal/constants"
)

const (
	PLAYER_COUNT_REGEX = `Players connected \((\d+)\):`
)

// Project Zomboid related commands
func pz(i discord.Interaction) (string, error) {
	subCommand := i.(discord.ApplicationCommandInteraction).SlashCommandInteractionData()

	switch *subCommand.SubCommandName {
	case "rcon": // /pz rcon <command>
		output, err := rcon.Run(os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PASSWORD"), subCommand.String("command"))
		if err != nil {
			return "", err
		}

		if output == "" {
			return "Successfully sent the RCON command", nil
		} else {
			return fmt.Sprintf("Successfully sent the RCON command and received the output `%s`", output), nil
		}
	case "start": // /pz start
		if err := projectzomboid.Start(); err != nil {
			return "", err
		}

		return constants.PZ_STATUS_ONLINE, nil
	case "status": // /pz status
		status, err := projectzomboid.Status()
		if err != nil {
			return "", err
		}

		if status {
			return constants.PZ_STATUS_ONLINE, nil
		} else {
			return constants.PZ_STATUS_OFFLINE, nil
		}
	case "stop": // /pz stop
		if err := projectzomboid.Stop(); err != nil {
			return "", err
		}

		return constants.PZ_STATUS_OFFLINE, nil
	}

	return "", errors.New("invalid pz subcommand")
}
