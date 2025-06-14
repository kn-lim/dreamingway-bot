package commands

import (
	"dario.cat/mergo"
	"github.com/disgoorg/disgo/discord"
)

type Command struct {
	Command discord.ApplicationCommandCreate
	Handler func(discord.Interaction) (string, error)
	// TODO: Implement options for commands
}

var (
	GlobalCommands = map[string]Command{
		"coinflip": {
			Command: discord.SlashCommandCreate{
				Name:        "coinflip",
				Description: "Flips a coin",
			},
			Handler: coinflip,
		},
		"ping": {
			Command: discord.SlashCommandCreate{
				Name:        "ping",
				Description: "Ping",
			},
			Handler: ping,
		},
		"roll": {
			Command: discord.SlashCommandCreate{
				Name:        "roll",
				Description: "Rolls a dice with modifiers",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "dice",
						Description: "Amount of dice to roll plus modifiers",
						Required:    true,
					},
				},
			},
			Handler: roll,
		},
	}

	Commands = map[string]Command{}
)

// GetAllCommands returns a slice of all commands, both global and non-global
func GetAllCommands() (map[string]Command, error) {
	allCommands := make(map[string]Command)
	if err := mergo.Merge(&allCommands, GlobalCommands); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&allCommands, Commands); err != nil {
		return nil, err
	}
	return allCommands, nil
}
