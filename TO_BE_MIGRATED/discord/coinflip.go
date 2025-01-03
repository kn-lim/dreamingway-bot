package discord

import (
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// coinflip simulates flipping a coin by returning either "Heads" or "Tails" randomly
func coinflip(i *discordgo.Interaction, opts ...Option) (string, error) { // nolint:unusedparam
	log.Println("/coinflip")

	// Create a new rand.Rand instance with a seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number (0 or 1)
	if r.Intn(2) == 0 {
		return "Heads", nil
	} else {
		return "Tails", nil
	}
}
