package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

const (
	ErrorMessage = "Error! Something went wrong! Check the logs for more details."
)

func handler(interaction discordgo.Interaction) error {
	// log.Printf("Received interaction: %+v\n", interaction)

	cmd, ok := discord.Commands[interaction.ApplicationCommandData().Name]
	if !ok {
		log.Printf("Error! Command does not exist: %s", interaction.ApplicationCommandData().Name)
		return discord.SendDeferredMessage(interaction.AppID, interaction.Token, ErrorMessage)
	}

	if cmd.Handler != nil {
		log.Printf("Running the handler of %s", interaction.ApplicationCommandData().Name)

		message, err := cmd.Handler()
		if err != nil {
			return discord.SendDeferredMessage(interaction.AppID, interaction.Token, ErrorMessage)
		}

		return discord.SendDeferredMessage(interaction.AppID, interaction.Token, message)
	} else if cmd.Options[interaction.ApplicationCommandData().Options[0].Name] != nil {
		log.Printf("Running the option handlers of %s", interaction.ApplicationCommandData().Name)

		message, err := cmd.Options[interaction.ApplicationCommandData().Options[0].Name]()
		if err != nil {
			return discord.SendDeferredMessage(interaction.AppID, interaction.Token, ErrorMessage)
		}

		return discord.SendDeferredMessage(interaction.AppID, interaction.Token, message)
	}

	return discord.SendDeferredMessage(interaction.AppID, interaction.Token, ErrorMessage)
}

func main() {
	lambda.Start(handler)
}
