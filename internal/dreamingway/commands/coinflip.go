package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/chattingway/gamble"
)

// coinflip returns a string of the result of a coin flip
func coinflip(i *discordgo.Interaction) (string, error) {
	return fmt.Sprintf("Flipped `%s`", gamble.CoinFlip()), nil
}
