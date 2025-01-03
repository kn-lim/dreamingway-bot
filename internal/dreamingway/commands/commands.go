package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(*discordgo.Interaction) (string, error)
	Options map[string]func(*discordgo.Interaction) (string, error)
}

var (
	// Map of Command structs
	Commands = map[string]Command{
		"coinflip": {
			Command: discordgo.ApplicationCommand{
				Name:        "coinflip",
				Description: "Flips a coin",
			},
			Handler: coinflip,
			Options: nil,
		},
		"foundry": {
			Command: discordgo.ApplicationCommand{
				Name:        "foundry",
				Description: "Manage FoundryVTT",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "status",
						Description: "Get the status of the FoundryVTT server",
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "start",
						Description: "Starts the FoundryVTT server",
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "stop",
						Description: "Stops the FoundryVTT server",
					},
				},
			},
			Handler: nil,
			Options: map[string]func(*discordgo.Interaction) (string, error){
				"status": foundryStatus,
				"start":  foundryStart,
				"stop":   foundryStop,
			},
		},
		"ping": {
			Command: discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping",
			},
			Handler: ping,
			Options: nil,
		},
		"roll": {
			Command: discordgo.ApplicationCommand{
				Name:        "roll",
				Description: "Rolls a dice with modifiers",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "dice",
						Description: "Amount of dice to roll plus modifiers",
						Required:    true,
					},
				},
			},
			Handler: roll,
			Options: nil,
		},
	}
)
