package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/v2/gamble"
)

// coinflip returns a string of the result of a coin flip
func coinflip(i discord.Interaction) (string, error) {
	return fmt.Sprintf("Flipped `%s`", gamble.CoinFlip()), nil
}
