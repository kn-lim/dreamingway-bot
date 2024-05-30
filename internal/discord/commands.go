package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(*discordgo.Interaction, ...Option) (string, error)
	Options map[string]func(*discordgo.Interaction, ...Option) (string, error)
}

var (
	Commands = map[string]Command{
		"coinflip": {
			Command: discordgo.ApplicationCommand{
				Name:        "Coinflip",
				Description: "Flip a coin",
			},
			Handler: coinflip,
			Options: nil,
		},
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
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "start",
						Description: "Starts the Pixelmon server",
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "stop",
						Description: "Stops the Pixelmon server",
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "say",
						Description: "Sends a message to the Pixelmon server",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "message",
								Description: "Message to send to the Pixelmon server",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "whitelist",
						Description: "Adds a user to the whitelist of the Pixelmon server",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "username",
								Description: "Minecraft username to whitelist",
								Required:    true,
							},
						},
					},
				},
			},
			Handler: nil,
			Options: map[string]func(*discordgo.Interaction, ...Option) (string, error){
				"status":    status,
				"start":     start,
				"stop":      stop,
				"say":       say,
				"whitelist": whitelist,
			},
		},
		"roll": {
			Command: discordgo.ApplicationCommand{
				Name:        "roll",
				Description: "Roll the dice",
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
