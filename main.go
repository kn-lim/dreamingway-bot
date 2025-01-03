package main

import (
	"flag"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/goccy/go-yaml"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

var (
	s *discordgo.Session

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
		utils.Logger.Errorw("Error reading YAML file",
			"error", err,
		)
		os.Exit(1)
	}

	// Unmarshal the YAML data into the Config struct
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		utils.Logger.Errorw("Error unmarshalling YAML",
			"error", err,
		)
		os.Exit(1)
	}

	// Create a new Discord session
	s, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		utils.Logger.Errorw("Error creating Discord session",
			"error", err,
		)
		os.Exit(1)
	}
	if err := s.Open(); err != nil {
		utils.Logger.Errorw("Error opening Discord session",
			"error", err,
		)
		os.Exit(1)
	}

	for _, server := range cfg.Servers {
		// Get all commands currently available for the bot
		curr_cmds, err := s.ApplicationCommands(cfg.AppID, server.GuildID)
		if err != nil {
			utils.Logger.Errorw("Error getting current commands",
				"error", err,
			)
			os.Exit(1)
		}

		// Delete all commands
		for _, cmd := range curr_cmds {
			if err := s.ApplicationCommandDelete(cfg.AppID, server.GuildID, cmd.ID); err != nil {
				utils.Logger.Errorw("Error deleting command",
					"command", cmd.Name,
					"error", err,
				)
				os.Exit(1)
			}
		}
		utils.Logger.Infow("Finished deleting all commands",
			"server", server.GuildID,
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
					if _, err := s.ApplicationCommandCreate(cfg.AppID, server.GuildID, &cmd); err != nil {
						utils.Logger.Errorw("Error uploading command",
							"command", cmd.Name,
							"server", server.GuildID,
							"error", err,
						)
						os.Exit(1)
					} else {
						utils.Logger.Infow("Successfully uploaded command",
							"command", cmd.Name,
							"server", server.GuildID,
						)
					}
					continue
				}
			}
		}
	}
}
