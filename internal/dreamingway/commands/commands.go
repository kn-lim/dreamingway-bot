package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(*discordgo.Interaction) (string, error)
	Options map[string]func(*discordgo.Interaction) (string, error)
}

var (
	Commands = map[string]Command{
		"coinflip": {
			Command: discordgo.ApplicationCommand{
				Name:        "coinflip",
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
