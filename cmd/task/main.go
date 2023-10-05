package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

func handler(interaction discordgo.Interaction) error {
	// log.Printf("Received interaction: %+v\n", interaction)

	return discord.SendDeferredMessage(interaction.AppID, interaction.Token, "This is from the task lambda!")
}

func main() {
	lambda.Start(handler)
}
