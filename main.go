package main

import (
	"flag"
	"fmt"
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
	GuildID string `yaml:"guild_id"`
	Token   string `yaml:"token"`
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
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Unmarshal the YAML data into the Config struct
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	// Create a new Discord session
	s, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("Error! Invalid bot parameters: %v", err)
	}

	if err := s.Open(); err != nil {
		log.Fatalf("Error! Cannot open the session: %v", err)
	}

	// Get all commands currently available for the bot
	curr_cmds, err := s.ApplicationCommands(cfg.AppID, cfg.GuildID)
	if err != nil {
		log.Fatal("Error! Couldn't get current commands!\n")
	}

	// Delete all commands
	for _, cmd := range curr_cmds {
		if err := s.ApplicationCommandDelete(cfg.AppID, cfg.GuildID, cmd.ID); err != nil {
			log.Fatalf("Error! Couldn't delete command [%s]: %v", cmd.Name, err)
		}
	}
	fmt.Println("Finished deleting all commands")

	// Add commands
	for i := range commands.Commands {
		cmd := commands.Commands[i].Command
		if _, err := s.ApplicationCommandCreate(cfg.AppID, cfg.GuildID, &cmd); err != nil {
			log.Fatalf("Error! Couldn't upload command [%s]: %v", cmd.Name, err)
		} else {
			fmt.Printf("Successfully uploaded command [%s]\n", cmd.Name)
		}
	}
}
