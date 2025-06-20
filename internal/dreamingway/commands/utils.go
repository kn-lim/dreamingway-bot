package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"

	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

// SyncGlobalCommands syncs the global commands with the provided commands
func SyncGlobalCommands(client rest.Rest, applicationID snowflake.ID, commands map[string]Command) error {
	// Get current global commands
	curr_cmds, err := client.GetGlobalCommands(applicationID, false)
	if err != nil {
		return err
	}

	// Delete all global commands
	var deletedCmds []string
	for _, cmd := range curr_cmds {
		if err := client.DeleteGlobalCommand(applicationID, cmd.ID()); err != nil {
			utils.Logger.Errorw("failed to delete global command",
				"command_id", cmd.ID(),
				"command_name", cmd.Name(),
			)
			return err
		}

		deletedCmds = append(deletedCmds, cmd.Name())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("deleted global commands",
			"commands", deletedCmds,
		)
	}

	var cmdSlice []discord.ApplicationCommandCreate
	for _, cmd := range commands {
		cmdSlice = append(cmdSlice, cmd.Command)
	}

	// Set new global commands
	if _, err := client.SetGlobalCommands(applicationID, cmdSlice); err != nil {
		utils.Logger.Errorw("failed to set global commands",
			"commands", cmdSlice,
		)
		return err
	}
	var newCmds []string
	for _, cmd := range cmdSlice {
		newCmds = append(newCmds, cmd.CommandName())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("synced global commands",
			"commands", newCmds,
		)
	}

	return nil
}

// SyncGuildCommands syncs the guild commands with the provided commands
func SyncGuildCommands(client rest.Rest, applicationID, guildID snowflake.ID, commands map[string]Command) error {
	// Get current guild commands
	curr_cmds, err := client.GetGuildCommands(applicationID, guildID, false)
	if err != nil {
		return err
	}

	// Delete all guild commands
	var deletedCmds []string
	for _, cmd := range curr_cmds {
		if err := client.DeleteGuildCommand(applicationID, guildID, cmd.ID()); err != nil {
			utils.Logger.Errorw("failed to delete guild command",
				"command_id", cmd.ID(),
				"command_name", cmd.Name(),
			)
			return err
		}

		deletedCmds = append(deletedCmds, cmd.Name())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("deleted guild commands",
			"guild_id", guildID,
			"commands", deletedCmds,
		)
	}

	var cmdSlice []discord.ApplicationCommandCreate
	for _, cmd := range commands {
		cmdSlice = append(cmdSlice, cmd.Command)
	}

	// Set new guild commands
	if _, err := client.SetGuildCommands(applicationID, guildID, cmdSlice); err != nil {
		utils.Logger.Errorw("failed to set guild commands",
			"guild_id", guildID,
			"commands", cmdSlice,
		)
		return err
	}
	var newCmds []string
	for _, cmd := range cmdSlice {
		newCmds = append(newCmds, cmd.CommandName())
	}
	if utils.Logger != nil {
		utils.Logger.Infow("synced guild commands",
			"guild_id", guildID,
			"commands", newCmds,
		)
	}

	return nil
}
