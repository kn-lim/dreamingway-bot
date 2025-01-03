package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/chattingway/gamble"
)

// roll returns a string of the result of a dice roll
func roll(i *discordgo.Interaction) (string, error) {
	input := i.ApplicationCommandData().Options[0].StringValue()
	output, value, err := gamble.Roll(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = **%d**", output, value), nil
}
