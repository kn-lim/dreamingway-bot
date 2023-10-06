package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ping(i *discordgo.Interaction) (string, error) {
	log.Println("/ping")

	return "Pong!", nil
}
