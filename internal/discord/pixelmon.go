package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

func status(*discordgo.Interaction) (discordgo.InteractionResponse, error) {
	// log.Println("status")

	result, online, err := mcstatus.GetMCStatus()
	if err != nil {
		log.Printf("Error! Couldn't get status: %s", err)
	}

	var status string
	if result {
		status = fmt.Sprintf("%v | Number of Online Players: %v", fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN")), online)
	} else {
		status = "Offline!"
	}

	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: status,
		},
	}, nil
}
