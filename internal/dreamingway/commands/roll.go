package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/gamble"
)

// roll returns a string of the result of a dice roll
func roll(i discord.Interaction) (string, error) {
	output, value, err := gamble.Roll(i.(discord.ApplicationCommandInteraction).Data.(discord.SlashCommandInteractionData).String("dice"))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Rolled `%s` = **%d**", output, value), nil
}
