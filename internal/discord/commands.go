package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(*discordgo.Interaction) (discordgo.InteractionResponse, error)
	Options map[string]func(*discordgo.Interaction) (discordgo.InteractionResponse, error)
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
			Handler: DeferredMessage,
			Options: map[string]func(*discordgo.Interaction) (discordgo.InteractionResponse, error){
				"status": status,
			},
		},
	}
)
