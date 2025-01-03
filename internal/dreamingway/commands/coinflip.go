package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/chattingway/gamble"
)

func coinflip(i *discordgo.Interaction) (string, error) {
	return gamble.CoinFlip(), nil
}
