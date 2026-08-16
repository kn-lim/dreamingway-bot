package commands

import (
	"github.com/disgoorg/disgo/discord"
)

// Command pairs a Discord application command with its handler.
type Command struct {
	Command discord.ApplicationCommandCreate
	Handler func(discord.Interaction) (string, error)
}

// Commands holds every command the bot can run.
// Each bot's config file decides which of these commands that bot publishes.
var Commands = map[string]Command{
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
			Description: "Rolls dice with modifiers",
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

	// counters
	"counter": {
		Command: discord.SlashCommandCreate{
			Name:        "counter",
			Description: "Manage counters",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "set",
					Description: "Set a counter to an exact value",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "Name of the counter",
							Required:    true,
						},
						discord.ApplicationCommandOptionInt{
							Name:        "value",
							Description: "New value for the counter",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "message",
					Description: "Set the reply message of a counter",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "Name of the counter",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "message",
							Description: "Reply message. Use {count} for the value",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "list",
					Description: "Show every counter in this server",
				},
			},
		},
		Handler: counterAdmin,
	},
	"test1": newCounter("test1", "Test #1 counter"),
	"test2": newCounter("test2", "Test #2 counter"),
}
