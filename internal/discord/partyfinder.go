package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func partyfinder(i *discordgo.Interaction, opts ...Option) (string, error) {
	log.Println("/partyfinder")

	return fmt.Sprintf("/partyfinder %v %v %v %v", i.ApplicationCommandData().Options[0].Value, i.ApplicationCommandData().Options[1].Value, i.ApplicationCommandData().Options[2].Value, i.ApplicationCommandData().Options[3].Value), nil
}
