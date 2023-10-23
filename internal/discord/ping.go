package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ping(i *discordgo.Interaction, opts ...Option) (string, error) {
	log.Println("/ping")

	return "Pong!", nil
}
