package main

import (
	"flag"
	"log"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

var (
	ConfigPath = flag.String("config", "config.yaml", "Path to the config file")
)

type Config struct {
	AppID   string `yaml:"app_id"`
	Token   string `yaml:"token"`
	Servers []struct {
		GuildID  string   `yaml:"guild_id"`
		Commands []string `yaml:"commands"`
	} `yaml:"servers"`
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
	// Read the YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		utils.Logger.Errorw("failed to read YAML file",
			"error", err,
		)
		os.Exit(1)
	}

	// Unmarshal the YAML data into the Config struct
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		utils.Logger.Errorw("failed to unmarshal YAML",
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
	if err := d.Client.Open(); err != nil {
		utils.Logger.Errorw("failed to open Discord session",
			"error", err,
		)
		os.Exit(1)
	}

	for _, server := range cfg.Servers {
		// Get server name from guild_id
		serverName, err := d.GetServerName(server.GuildID)
		if err != nil {
			utils.Logger.Errorw("failed to get server name",
				"guild_id", server.GuildID,
				"error", err,
			)
			os.Exit(1)
		}

		// Get all commands currently available for the bot
		curr_cmds, err := d.Client.ApplicationCommands(cfg.AppID, server.GuildID)
		if err != nil {
			utils.Logger.Errorw("failed to get current commands",
				"error", err,
			)
			os.Exit(1)
		}

		// Delete all commands
		for _, cmd := range curr_cmds {
			if err := d.Client.ApplicationCommandDelete(cfg.AppID, server.GuildID, cmd.ID); err != nil {
				utils.Logger.Errorw("failed to delete command",
					"command", cmd.Name,
					"error", err,
				)
				os.Exit(1)
			}
		}
		utils.Logger.Infow("finished deleting all commands",
			"server", serverName,
			"guild_id", server.GuildID,
		)

		// Add commands
		var cmds []string
		if len(server.Commands) == 0 {
			for i := range commands.Commands {
				cmds = append(cmds, commands.Commands[i].Command.Name)
			}
		} else {
			cmds = server.Commands
		}
		for i := range commands.Commands {
			for _, cmd := range cmds {
				if commands.Commands[i].Command.Name == cmd {
					cmd := commands.Commands[i].Command
					if _, err := d.Client.ApplicationCommandCreate(cfg.AppID, server.GuildID, &cmd); err != nil {
						utils.Logger.Errorw("failed to upload command",
							"command", cmd.Name,
							"server", serverName,
							"guild_id", server.GuildID,
							"error", err,
						)
						os.Exit(1)
					} else {
						utils.Logger.Infow("successfully uploaded command",
							"command", cmd.Name,
							"server", serverName,
							"guild_id", server.GuildID,
						)
					}
					continue
				}
			}
		}
	}
}
