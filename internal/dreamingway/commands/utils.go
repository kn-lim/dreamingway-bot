package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

// SyncGlobalCommands syncs the global commands with the provided commands
func SyncGlobalCommands(client rest.Rest, applicationID snowflake.ID, commands []discord.ApplicationCommandCreate) error {
	// Get current global commands
	curr_cmds, err := client.GetGlobalCommands(applicationID, false)
	if err != nil {
		return err
	}

	// Delete all global commands
	var deletedCmds []string
	for _, cmd := range curr_cmds {
		if err := client.DeleteGlobalCommand(applicationID, cmd.ApplicationID()); err != nil {
			return err
		}

		deletedCmds = append(deletedCmds, cmd.Name())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("deleted global commands",
			"application_id", applicationID,
			"commands", deletedCmds,
		)
	}

	// Set new global commands
	if _, err := client.SetGlobalCommands(applicationID, commands); err != nil {
		return err
	}
	var newCmds []string
	for _, cmd := range commands {
		newCmds = append(newCmds, cmd.CommandName())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("synced global commands",
			"application_id", applicationID,
			"commands", newCmds,
		)
	}

	return nil
}

// SyncGuildCommands syncs the guild commands with the provided commands
func SyncGuildCommands(client rest.Rest, applicationID, guildID snowflake.ID, commands []discord.ApplicationCommandCreate) error {
	// Get current guild commands
	curr_cmds, err := client.GetGuildCommands(applicationID, guildID, false)
	if err != nil {
		return err
	}

	// Delete all guild commands
	var deletedCmds []string
	for _, cmd := range curr_cmds {
		if err := client.DeleteGuildCommand(applicationID, guildID, cmd.ApplicationID()); err != nil {
			return err
		}

		deletedCmds = append(deletedCmds, cmd.Name())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("deleted guild commands",
			"application_id", applicationID,
			"guild_id", guildID,
			"commands", deletedCmds,
		)
	}

	// Set new guild commands
	if _, err := client.SetGuildCommands(applicationID, guildID, commands); err != nil {
		return err
	}
	var newCmds []string
	for _, cmd := range commands {
		newCmds = append(newCmds, cmd.CommandName())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("synced guild commands",
			"application_id", applicationID,
			"guild_id", guildID,
			"commands", newCmds,
		)
	}

	return nil
}
