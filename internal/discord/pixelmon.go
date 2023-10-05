package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/pixelmon"
)

func status(interaction *discordgo.Interaction) (discordgo.InteractionResponse, error) {
	log.Println("status")

	result, online, err := pixelmon.GetStatus()
	if err != nil {
		log.Printf("Error! Couldn't get status: %s", err)
	}

	var status string
	serverURL := fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
	if result {
		log.Printf("%v is online", serverURL)
		status = fmt.Sprintf("%v | Number of Online Players: %v", serverURL, online)
	} else {
		log.Printf("%v is offline", serverURL)
		status = fmt.Sprintf("%v | Currently Offline", serverURL)
	}

	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: status,
		},
	}, nil
}
