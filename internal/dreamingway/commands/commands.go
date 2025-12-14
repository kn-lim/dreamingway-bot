package commands

import (
	"dario.cat/mergo"
	"github.com/disgoorg/disgo/discord"
)

type Command struct {
	Command discord.ApplicationCommandCreate
	Handler func(discord.Interaction) (string, error)
}

var (
	GlobalCommands = map[string]Command{
		// healthcheck
		"ping": {
			Command: discord.SlashCommandCreate{
				Name:        "ping",
				Description: "Ping",
			},
			Handler: ping,
		},

		// gamble
		"coinflip": {
			Command: discord.SlashCommandCreate{
				Name:        "coinflip",
				Description: "Flip a coin",
			},
			Handler: coinflip,
		},
		"roll": {
			Command: discord.SlashCommandCreate{
				Name:        "roll",
				Description: "Roll a dice with modifiers",
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

		// games
		"pz": {
			Command: discord.SlashCommandCreate{
				Name:        "pz",
				Description: "Project Zomboid related commands",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionSubCommand{
						Name:        "rcon",
						Description: "Send an RCON command to the Project Zomboid server",
						Options: []discord.ApplicationCommandOption{
							discord.ApplicationCommandOptionString{
								Name:        "command",
								Description: "RCON command",
								Required:    true,
							},
						},
					},
					discord.ApplicationCommandOptionSubCommand{
						Name:        "start",
						Description: "Start the Project Zomboid server",
					},
					discord.ApplicationCommandOptionSubCommand{
						Name:        "status",
						Description: "Get the status of the Project Zomboid server",
					},
					discord.ApplicationCommandOptionSubCommand{
						Name:        "stop",
						Description: "Stop the Project Zomboid server",
					},
				},
			},
			Handler: pz,
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
