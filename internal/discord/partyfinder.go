package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func partyfinder(i *discordgo.Interaction, opts ...Option) (string, error) {
	log.Println("/partyfinder")

	// Create actions row
	actionsRow := &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			// Create buttons
			discordgo.Button{
				Label:    "Tank",
				Style:    discordgo.PrimaryButton,
				CustomID: "tank",
			},
			discordgo.Button{
				Label:    "Healer",
				Style:    discordgo.SuccessButton,
				CustomID: "healer",
			},
			discordgo.Button{
				Label:    "DPS",
				Style:    discordgo.DangerButton,
				CustomID: "dps",
			},
			discordgo.Button{
				Label:    "🗑️",
				Style:    discordgo.SecondaryButton,
				CustomID: "delete",
			},
		},
	}

	if err := SendDeferredMessageWithComponents(i.AppID, i.Token, "/partyfinder", WithActionsRow(actionsRow)); err != nil {
		return "", err
	}

	return "", nil
}
