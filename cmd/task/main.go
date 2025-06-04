package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/disgoorg/disgo/discord"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

func handler(ctx context.Context, rawInteraction json.RawMessage) error {
	// Initialize logger
	var err error
	utils.Logger, err = utils.NewLogger(true)
	if err != nil {
		log.Printf("couldn't initialize logger: %v", err)
		return err
	}

	// Get discord interaction
	interaction, err := discord.UnmarshalInteraction(rawInteraction)
	if err != nil {
		utils.Logger.Errorw("couldn't unmarshal interaction",
			"error", err,
			"interaction", string(rawInteraction),
		)
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

	// Get command
	cmds, err := commands.GetAllCommands()
	if err != nil {
		utils.Logger.Errorw("couldn't get commands",
			"error", err,
		)
		return dreamingwayBot.SendDeferredMessage(
			interaction.ApplicationID().String(),
			interaction.Token(),
			"**Error**! Could not get commands.",
		)
	}
	cmd, ok := cmds[interaction.(discord.ApplicationCommandInteraction).Data.CommandName()]
	if !ok {
		utils.Logger.Errorw("command does not exist",
			"command", interaction.(discord.ApplicationCommandInteraction).Data.CommandName(),
			"username", interaction.User().Username,
			"guild_id", interaction.GuildID().String(),
		)
		return dreamingwayBot.SendDeferredMessage(
			interaction.ApplicationID().String(),
			interaction.Token(),
			fmt.Sprintf("**Error**! Command `/`%s does not exist.", interaction.(discord.ApplicationCommandInteraction).Data.CommandName()),
		)
	}

	// Run command handler
	var msg string
	if cmd.Handler != nil {
		utils.Logger.Infow("running command handler",
			"command", interaction.(discord.ApplicationCommandInteraction).Data.CommandName(),
			"username", interaction.User().Username,
			"guild_id", interaction.GuildID().String(),
		)

		msg, err = cmd.Handler(interaction)
		if err != nil {
			utils.Logger.Errorw("error running command handler",
				"command", interaction.(discord.ApplicationCommandInteraction).Data.CommandName(),
				"username", interaction.User().Username,
				"guild_id", interaction.GuildID().String(),
				"error", err,
			)
			return dreamingwayBot.SendDeferredMessage(
				interaction.ApplicationID().String(),
				interaction.Token(),
				fmt.Sprintf("**Error**! /%s handler failed: `%s`", interaction.(discord.ApplicationCommandInteraction).Data.CommandName(), err),
			)
		}
		// TODO: Implement options for commands
	}

	if msg == "" {
		utils.Logger.Errorw("got empty message",
			"command", interaction.(discord.ApplicationCommandInteraction).Data.CommandName(),
			"username", interaction.User().Username,
			"guild_id", interaction.GuildID().String(),
		)
		return dreamingwayBot.SendDeferredMessage(
			interaction.ApplicationID().String(),
			interaction.Token(),
			fmt.Sprintf("**Error**! Got empty message for /%s.", interaction.(discord.ApplicationCommandInteraction).Data.CommandName()),
		)
	}

	return dreamingwayBot.SendDeferredMessage(
		interaction.ApplicationID().String(),
		interaction.Token(),
		msg,
	)
}

func main() {
	lambda.Start(handler)
}
