package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(*discordgo.Interaction) (string, error)
	Options map[string]func(*discordgo.Interaction) (string, error)
}

var (
	Commands = map[string]Command{
		"ping": {
			Command: discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Pong!",
			},
			Handler: ping,
			Options: nil,
		},
		"pixelmon": {
			Command: discordgo.ApplicationCommand{
				Name:        "pixelmon",
				Description: "Pixelmon command",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "status",
						Description: "Get the status of the Pixelmon server",
					},
				},
			},
			Handler: nil,
			Options: map[string]func(*discordgo.Interaction) (string, error){
				"status": status,
			},
		},
	}
)
