package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

var (
	s *discordgo.Session

	AppID   = flag.String("appid", "", "Discord Application ID")
	GuildID = flag.String("guildid", "", "Discord Guild ID (Optional)")
	Token   = flag.String("token", "", "Discord Bot Token")
)

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *Token)
	if err != nil {
		log.Fatalf("Error! Invalid bot parameters: %v", err)
	}
}

func main() {
	if err := s.Open(); err != nil {
		log.Fatalf("Error! Cannot open the session: %v", err)
	}

	// Get all commands currently available for the bot
	curr_cmds, err := s.ApplicationCommands(*AppID, *GuildID)
	if err != nil {
		log.Fatal("Error! Couldn't get current commands!\n")
	}

	// Delete all commands
	for _, cmd := range curr_cmds {
		if err := s.ApplicationCommandDelete(*AppID, *GuildID, cmd.ID); err != nil {
			log.Fatalf("Error! Couldn't delete command [%s]: %v", cmd.Name, err)
		}
	}
	fmt.Println("Finished deleting all commands")

	// Add commands
	for i := range discord.Commands {
		cmd := discord.Commands[i].Command
		if _, err := s.ApplicationCommandCreate(*AppID, *GuildID, &cmd); err != nil {
			log.Fatalf("Error! Couldn't upload command [%s]: %v", cmd.Name, err)
		} else {
			fmt.Printf("Successfully uploaded command [%s]\n", cmd.Name)
		}
	}
}
