package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

var (
	ConfigPath = flag.String("config", "config.json", "Path to the config file")
)

type Config struct {
	AppID          string `json:"app_id"`
	Token          string `json:"token"`
	GlobalCommands string `json:"global_commands"`
	Servers        []struct {
		GuildID  string   `json:"guild_id"`
		Commands []string `json:"commands"`
	} `json:"servers"`
}

func init() {
	flag.Parse()

	// Initialize logger
	var err error
	utils.Logger, err = utils.NewLogger()
	if err != nil {
		log.Printf("couldn't initialize logger: %v", err)
		os.Exit(1)
	}
}

func main() {
	// Read the config file
	var configFilePath string
	if ConfigPath != nil {
		_, err := os.Stat(*ConfigPath)
		if err != nil || os.IsNotExist(err) {
			utils.Logger.Errorw("config file not found",
				"error", err,
			)
			os.Exit(1)
		}

		configFilePath = *ConfigPath
	} else {
		configFilePath = "config.json"
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		utils.Logger.Errorw("failed to read config file",
			"error", err,
		)
		os.Exit(1)
	}

	// Unmarshal the JSON data into the Config struct
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		utils.Logger.Errorw("failed to unmarshal JSON",
			"error", err,
		)
		os.Exit(1)
	}

	// Create a new Discord session
	d, err := dreamingway.NewDreamingway(cfg.Token)
	if err != nil {
		utils.Logger.Errorw("failed to create Discord session",
			"error", err,
		)
		os.Exit(1)
	}
	// if err := d.Client.Open(); err != nil {
	// 	utils.Logger.Errorw("failed to open Discord session",
	// 		"error", err,
	// 	)
	// 	os.Exit(1)
	// }

	applicationID, err := snowflake.Parse(cfg.AppID)
	if err != nil {
		utils.Logger.Errorw("failed to parse application ID",
			"application_id", cfg.AppID,
			"error", err,
		)
		os.Exit(1)
	}

	// Sync global commands
	if err := commands.SyncGlobalCommands(d.Client.Rest(), applicationID, commands.GlobalCommands); err != nil {
		utils.Logger.Errorw("failed to sync global commands",
			"error", err,
		)
		os.Exit(1)
	}

	for _, server := range cfg.Servers {
		// Get server name from guild_id
		// serverName, err := d.GetServerName(server.GuildID)
		// if err != nil {
		// 	utils.Logger.Errorw("failed to get server name",
		// 		"guild_id", server.GuildID,
		// 		"error", err,
		// 	)
		// 	os.Exit(1)
		// }
		snowflakeID, err := snowflake.Parse(server.GuildID)
		if err != nil {
			utils.Logger.Errorw("failed to parse guild ID",
				"guild_id", server.GuildID,
				"error", err,
			)
			os.Exit(1)
		}
		// guild, err := d.Client.Rest().GetGuild(snowflakeID, false)
		// if err != nil {
		// 	utils.Logger.Errorw("failed to get guild",
		// 		"guild_id", server.GuildID,
		// 		"error", err,
		// 	)
		// 	os.Exit(1)
		// }

		// Check for server specific commands
		if len(server.Commands) == 0 {
			continue
		}

		// Sync guild commands
		var cmds []discord.ApplicationCommandCreate
		for i := range commands.Commands {
			for _, cmd := range server.Commands {
				if commands.Commands[i].CommandName() == cmd {
					cmds = append(cmds, commands.Commands[i])
				}
			}
		}
		if err := commands.SyncGuildCommands(d.Client.Rest(), applicationID, snowflakeID, cmds); err != nil {
			utils.Logger.Errorw("failed to sync guild commands",
				"guild_id", server.GuildID,
				"error", err,
			)
			os.Exit(1)
		}

		// Get all commands currently available for the bot
		// curr_cmds, err := d.Client.ApplicationCommands(cfg.AppID, server.GuildID)
		// if err != nil {
		// 	utils.Logger.Errorw("failed to get current commands",
		// 		"error", err,
		// 	)
		// 	os.Exit(1)
		// }

		// Delete all commands
		// for _, cmd := range curr_cmds {
		// 	if err := d.Client.ApplicationCommandDelete(cfg.AppID, server.GuildID, cmd.ID); err != nil {
		// 		utils.Logger.Errorw("failed to delete command",
		// 			"command", cmd.Name,
		// 			"error", err,
		// 		)
		// 		os.Exit(1)
		// 	}
		// }
		// utils.Logger.Infow("finished deleting all commands",
		// 	"server", guild.Name,
		// 	"guild_id", server.GuildID,
		// )

		// Add commands
		// var cmds []string
		// if len(server.Commands) == 0 {
		// 	for i := range commands.Commands {
		// 		cmds = append(cmds, commands.Commands[i].Command.Name)
		// 	}
		// } else {
		// 	cmds = server.Commands
		// }
		// for i := range commands.Commands {
		// 	for _, cmd := range cmds {
		// 		if commands.Commands[i].Command.Name == cmd {
		// 			cmd := commands.Commands[i].Command
		// 			if _, err := d.Client.ApplicationCommandCreate(cfg.AppID, server.GuildID, &cmd); err != nil {
		// 				utils.Logger.Errorw("failed to upload command",
		// 					"command", cmd.Name,
		// 					"server", serverName,
		// 					"guild_id", server.GuildID,
		// 					"error", err,
		// 				)
		// 				os.Exit(1)
		// 			} else {
		// 				utils.Logger.Infow("successfully uploaded command",
		// 					"command", cmd.Name,
		// 					"server", serverName,
		// 					"guild_id", server.GuildID,
		// 				)
		// 			}
		// 			continue
		// 		}
		// 	}
		// }
	}
}
