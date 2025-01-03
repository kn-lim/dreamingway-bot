package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/TO_BE_MIGRATED/dice"
)

func roll(i *discordgo.Interaction, opts ...Option) (string, error) {
	log.Println("/roll")

	// Get the dice value
	diceString := i.ApplicationCommandData().Options[0].StringValue()
	output, value, err := dice.Roll(diceString)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = **%d**", output, value), nil
}
