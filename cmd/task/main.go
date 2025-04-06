package main

import (
	"fmt"
	"log"
	"os"

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

	// Create a new Discord session
	dreamingwayBot, err := dreamingway.NewDreamingway(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		utils.Logger.Errorw("couldn't create a new Discord session",
			"error", err,
		)
		return err
	}

	// Get username
	username := dreamingway.GetUsername(interaction)

	// Get server name
	serverName, err := dreamingwayBot.GetServerName(interaction.GuildID)
	if err != nil {
		utils.Logger.Errorw("failed to get server name",
			"error", err,
			"guild_id", interaction.GuildID,
		)
		return err
	}

	// Get command
	cmd, ok := commands.Commands[interaction.ApplicationCommandData().Name]
	if !ok {
		utils.Logger.Errorw("command does not exist",
			"command", interaction.ApplicationCommandData().Name,
			"username", username,
			"server", serverName,
		)
		return dreamingwayBot.SendDeferredMessage(interaction.AppID, interaction.Token, "**Error**! Command does not exist.")
	}

	// Run command handler
	var msg string
	if cmd.Handler != nil {
		utils.Logger.Infow("running command handler",
			"command", interaction.ApplicationCommandData().Name,
			"username", username,
			"server", serverName,
		)

		msg, err = cmd.Handler(&interaction)
		if err != nil {
			utils.Logger.Errorw("error running command handler",
				"command", interaction.ApplicationCommandData().Name,
				"username", username,
				"server", serverName,
				"error", err,
			)
			return dreamingwayBot.SendDeferredMessage(interaction.AppID, interaction.Token, fmt.Sprintf("**Error**! /%s handler failed: `%s`", interaction.ApplicationCommandData().Name, err))
		}
	} else if cmd.Options[interaction.ApplicationCommandData().Options[0].Name] != nil {
		utils.Logger.Infow("running option handler",
			"command", interaction.ApplicationCommandData().Name,
			"username", username,
			"server", serverName,
		)

		msg, err = cmd.Options[interaction.ApplicationCommandData().Options[0].Name](&interaction)
		if err != nil {
			utils.Logger.Errorw("error running option handler",
				"command", interaction.ApplicationCommandData().Name,
				"username", username,
				"server", serverName,
				"error", err,
			)
			return dreamingwayBot.SendDeferredMessage(interaction.AppID, interaction.Token, fmt.Sprintf("**Error**! /%s option handler failed: `%s`", interaction.ApplicationCommandData().Name, err))
		}
	}

	if msg == "" {
		utils.Logger.Errorw("got empty message",
			"command", interaction.ApplicationCommandData().Name,
			"username", username,
			"server", serverName,
		)
		return dreamingwayBot.SendDeferredMessage(interaction.AppID, interaction.Token, fmt.Sprintf("**Error**! Got empty message for /%s.", interaction.ApplicationCommandData().Name))
	}

	return dreamingwayBot.SendDeferredMessage(interaction.AppID, interaction.Token, msg)
}

func main() {
	lambda.Start(handler)
}
