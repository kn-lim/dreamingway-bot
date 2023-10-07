package discord

import (
	"fmt"
	"log"

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

	if result {
		// log.Printf("%v is online", serverURL)
		return fmt.Sprintf(":green_circle:   %s | Online | Number of Online Players: %v", pixelmon.ServerURL, online), nil
	} else {
		// log.Printf("%v is offline", serverURL)
		return fmt.Sprintf(":red_circle:   %s | Offline", pixelmon.ServerURL), nil
	}
}

func start(i *discordgo.Interaction) (string, error) {
	log.Println("/pixelmon start")

	if err := SendDeferredMessage(i.AppID, i.Token, fmt.Sprintf(":green_square:   %s | Starting the Pixelmon server", pixelmon.ServerURL)); err != nil {
		return "", err
	}

	if err := pixelmon.StartInstance(); err != nil {
		return "", err
	}

	if err := pixelmon.StartService(); err != nil {
		return "", err
	}

	return fmt.Sprintf(":green_circle:   %s | Online", pixelmon.ServerURL), nil
}

func stop(i *discordgo.Interaction) (string, error) {
	log.Println("/pixelmon stop")

	if err := SendDeferredMessage(i.AppID, i.Token, fmt.Sprintf(":red_square:   %s | Stopping the Pixelmon server", pixelmon.ServerURL)); err != nil {
		return "", err
	}

	if err := pixelmon.StopInstance(); err != nil {
		return "", err
	}

	if err := pixelmon.StopService(); err != nil {
		return "", err
	}

	return fmt.Sprintf(":red_circle:   %s | Offline", pixelmon.ServerURL), nil
}

func say(i *discordgo.Interaction) (string, error) {
	return "", nil
}

func whitelist(i *discordgo.Interaction) (string, error) {
	return "", nil
}
