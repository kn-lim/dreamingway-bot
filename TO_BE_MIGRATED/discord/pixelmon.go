package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/TO_BE_MIGRATED/mcstatus"
	"github.com/kn-lim/dreamingway-bot/TO_BE_MIGRATED/pixelmon"
)

var (
	ServerURL = fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
)

func status(i *discordgo.Interaction, opts ...Option) (string, error) {
	log.Println("/pixelmon status")

	// Defaults
	config := &options{
		url: mcstatus.URL,
	}
	for _, opt := range opts {
		opt(config)
	}

	isOnline, playerCount, err := pixelmon.GetStatus(ServerURL, pixelmon.WithURL(config.url))
	if err != nil {
		log.Printf("Error! Couldn't get status: %s", err)
		return "", err
	}

	if isOnline {
		// log.Printf("%v is online", serverURL)
		return fmt.Sprintf(":green_circle:   %s | Online | Number of Online Players: %v", ServerURL, playerCount), nil
	} else {
		// log.Printf("%v is offline", serverURL)
		return fmt.Sprintf(":red_circle:   %s | Offline", ServerURL), nil
	}
}

// func start(i *discordgo.Interaction, opts ...Option) (string, error) {
// 	log.Println("/pixelmon start")

// 	// Check if user has correct role
// 	if !CheckRole(i.Member.Roles, os.Getenv("PIXELMON_ROLE_ID")) {
// 		return ErrMissingRole, nil
// 	}

// 	// Check if service is already running
// 	status, _, err := pixelmon.GetStatus(ServerURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	if status {
// 		return fmt.Sprintf(":green_circle:   %s | Online", ServerURL), nil
// 	}

// 	if err := SendDeferredMessage(i.AppID, i.Token, fmt.Sprintf(":green_square:   %s | Starting the Pixelmon server", ServerURL)); err != nil {
// 		return "", err
// 	}

// 	if err := pixelmon.StartInstance(os.Getenv("PIXELMON_INSTANCE_ID")); err != nil {
// 		return "", err
// 	}

// 	if err := pixelmon.StartService(os.Getenv("PIXELMON_INSTANCE_ID"), os.Getenv("PIXELMON_HOSTED_ZONE_ID"), ServerURL); err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf(":green_circle:   %s | Online", ServerURL), nil
// }

// func stop(i *discordgo.Interaction, opts ...Option) (string, error) {
// 	log.Println("/pixelmon stop")

// 	// Check if user has correct role
// 	if !CheckRole(i.Member.Roles, os.Getenv("PIXELMON_ROLE_ID")) {
// 		return ErrMissingRole, nil
// 	}

// 	// Check if service is already stopped
// 	status, _, err := pixelmon.GetStatus(ServerURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !status {
// 		return fmt.Sprintf(":red_circle:   %s | Offline", ServerURL), nil
// 	}

// 	if err := SendDeferredMessage(i.AppID, i.Token, fmt.Sprintf(":red_square:   %s | Stopping the Pixelmon server", ServerURL)); err != nil {
// 		return "", err
// 	}

// 	if err := pixelmon.StopService(os.Getenv("PIXELMON_INSTANCE_ID"), os.Getenv("PIXELMON_HOSTED_ZONE_ID"), ServerURL, os.Getenv("PIXELMON_RCON_PASSWORD")); err != nil {
// 		return "", err
// 	}

// 	if err := pixelmon.StopInstance(os.Getenv("PIXELMON_INSTANCE_ID")); err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf(":red_circle:   %s | Offline", ServerURL), nil
// }

// func say(i *discordgo.Interaction, opts ...Option) (string, error) {
// 	log.Printf("/pixelmon say")

// 	// Check if service is already stopped
// 	status, _, err := pixelmon.GetStatus(ServerURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !status {
// 		return fmt.Sprintf(":red_circle:   %s | Offline", ServerURL), nil
// 	}

// 	// Say message to server
// 	message := i.ApplicationCommandData().Options[0].Options[0].StringValue()
// 	if err := pixelmon.SayMessage(os.Getenv("PIXELMON_INSTANCE_ID"), os.Getenv("PIXELMON_RCON_PASSWORD"), i.Member.Nick, message); err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("Sent command to say `%s`.", message), nil
// }

// func whitelist(i *discordgo.Interaction, opts ...Option) (string, error) {
// 	log.Printf("/pixelmon whitelist")

// 	// Check if user has correct role
// 	if !CheckRole(i.Member.Roles, os.Getenv("PIXELMON_ROLE_ID")) {
// 		return ErrMissingRole, nil
// 	}

// 	// Check if service is already stopped
// 	status, _, err := pixelmon.GetStatus(ServerURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !status {
// 		return fmt.Sprintf(":red_circle:   %s | Offline", ServerURL), nil
// 	}

// 	// Add username to whitelist
// 	username := i.ApplicationCommandData().Options[0].Options[0].StringValue()
// 	if err := pixelmon.AddToWhitelist(os.Getenv("PIXELMON_INSTANCE_ID"), os.Getenv("PIXELMON_RCON_PASSWORD"), username); err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("Sent command to whitelist `%s`.", username), nil
// }
