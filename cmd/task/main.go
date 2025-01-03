package main

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

func handler(interaction discordgo.Interaction) error {
	// Initialize logger
	var err error
	utils.Logger, err = utils.NewLogger()
	if err != nil {
		log.Printf("couldn't initialize logger: %v", err)
		return err
	}

	// Get command
	cmd, ok := commands.Commands[interaction.ApplicationCommandData().Name]
	if !ok {
		utils.Logger.Errorw("command does not exist",
			"command", interaction.ApplicationCommandData().Name,
		)
		errorMsg := "Error! Command does not exist."
		return dreamingway.SendDeferredMessage(interaction.AppID, interaction.Token, errorMsg)
	}

	// Run command handler
	var msg string
	if cmd.Handler != nil {
		utils.Logger.Infow("running command handler",
			"command", interaction.ApplicationCommandData().Name,
		)

		msg, err = cmd.Handler(&interaction)
		if err != nil {
			utils.Logger.Errorw("error running command handler",
				"command", interaction.ApplicationCommandData().Name,
				"error", err,
			)
			return dreamingway.SendDeferredMessage(interaction.AppID, interaction.Token, "Error! Something went wrong.")
		}
	} else if cmd.Options[interaction.ApplicationCommandData().Options[0].Name] != nil {
		utils.Logger.Infow("running option handler",
			"command", interaction.ApplicationCommandData().Name,
		)

		msg, err = cmd.Options[interaction.ApplicationCommandData().Options[0].Name](&interaction)
		if err != nil {
			utils.Logger.Errorw("error running option handler",
				"command", interaction.ApplicationCommandData().Name,
				"error", err,
			)
			return dreamingway.SendDeferredMessage(interaction.AppID, interaction.Token, "Error! Something went wrong.")
		}
	}

	if msg == "" {
		utils.Logger.Errorw("got empty message",
			"command", interaction.ApplicationCommandData().Name,
		)
		return errors.New("empty message")
	}

	return dreamingway.SendDeferredMessage(interaction.AppID, interaction.Token, msg)
}

func main() {
	lambda.Start(handler)
}
