package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

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
		snowflakeID, err := snowflake.Parse(server.GuildID)
		if err != nil {
			utils.Logger.Errorw("failed to parse guild ID",
				"guild_id", server.GuildID,
				"error", err,
			)
			os.Exit(1)
		}
		// Check for server specific commands
		if len(server.Commands) == 0 {
			continue
		}

		// Sync guild commands
		cmds := make(map[string]commands.Command)
		for key, serverCmd := range commands.Commands {
			for _, cmd := range server.Commands {
				if serverCmd.Command.CommandName() == cmd {
					cmds[key] = serverCmd
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
	}
}
