package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/pixelmon"
)

func status(i *discordgo.Interaction) (string, error) {
	log.Println("/pixelmon status")

	result, online, err := pixelmon.GetStatus()
	if err != nil {
		log.Printf("Error! Couldn't get status: %s", err)
		return "", err
	}

	serverURL := fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
	if result {
		// log.Printf("%v is online", serverURL)
		return fmt.Sprintf(":green_circle:   %s | Number of Online Players: %v", serverURL, online), nil
	} else {
		// log.Printf("%v is offline", serverURL)
		return fmt.Sprintf(":red_circle:   %s | Currently Offline", serverURL), nil
	}
}

func start(i *discordgo.Interaction) (string, error) {
	log.Println("/pixelmon start")

	if err := SendDeferredMessage(i.AppID, i.Token, "test message 1"); err != nil {
		log.Println("Error with sending 1st deferred message")
		return "", err
	}

	return "test message 2!!!", nil
}
