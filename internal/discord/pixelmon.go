package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/kn-lim/dreamingway-bot/internal/pixelmon"
)

func status() (string, error) {
	log.Println("/pixelmon status")

	result, online, err := pixelmon.GetStatus()
	if err != nil {
		log.Printf("Error! Couldn't get status: %s", err)
	}

	serverURL := fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
	if result {
		log.Printf("%v is online", serverURL)
		return fmt.Sprintf("%v | Number of Online Players: %v", serverURL, online), nil
	} else {
		log.Printf("%v is offline", serverURL)
		return fmt.Sprintf("%v | Currently Offline", serverURL), nil
	}
}
