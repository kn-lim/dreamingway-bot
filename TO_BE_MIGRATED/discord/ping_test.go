package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestPing(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		i := &discordgo.Interaction{}
		r, err := ping(i)

		if err != nil {
			t.Fatalf("discord.ping() err = %v, want nil", err)
		}

		if r != "Pong!" {
			t.Fatalf("discord.ping() got \"%v\", want \"Pong!\"", r)
		}
	})
}
