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

	for i := range discord.Commands {
		cmd := discord.Commands[i].Command
		if _, err := s.ApplicationCommandCreate(*AppID, *GuildID, &cmd); err != nil {
			log.Fatalf("Error! Couldn't upload command [%s]: %v", cmd.Name, err)
		} else {
			fmt.Printf("Successfully uploaded command [%s]\n", cmd.Name)
		}
	}
}
