package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/v2/counter"
)

// counterCooldown is the time between two increments of one counter in one guild.
// Change this value to change the cooldown.
const counterCooldown = 1 * time.Minute

// countPlaceholder is the text that a reply template replaces with the count.
const countPlaceholder = "{count}"

// newCounter returns a command that increments the counter with the same name.
// Every counter command uses the same handler.
func newCounter(name, description string) Command {
	return Command{
		Command: discord.SlashCommandCreate{
			Name:        name,
			Description: description,
		},
		Handler: incrementCounter,
	}
}

// incrementCounter adds one to the counter named after the command.
func incrementCounter(i discord.Interaction) (string, error) {
	guildID := i.GuildID()
	if guildID == nil {
		return "", errors.New("command must be sent from a server")
	}

	name := i.(discord.ApplicationCommandInteraction).Data.CommandName()

	c, err := counter.Increment(
		context.TODO(),
		os.Getenv("COUNTER_TABLE_NAME"),
		os.Getenv("AWS_REGION"),
		guildID.String(),
		name,
		counterCooldown,
	)
	if err != nil {
		var cooldown *counter.CooldownError
		if errors.As(err, &cooldown) {
			return fmt.Sprintf("`/%s` is on cooldown. Try again in %s.",
				name, cooldown.Remaining.Round(time.Second)), nil
		}

		return "", err
	}

	return renderCounter(c), nil
}

// renderCounter builds the reply for a counter.
// If the template has the placeholder, renderCounter replaces it with the count.
// If the template has no placeholder, renderCounter adds the count at the end.
// If the counter has no template, renderCounter returns the name and the count.
func renderCounter(c counter.Counter) string {
	if c.Message == "" {
		return fmt.Sprintf("%s: %d", c.Name, c.Count)
	}

	count := strconv.FormatInt(c.Count, 10)
	if strings.Contains(c.Message, countPlaceholder) {
		return strings.ReplaceAll(c.Message, countPlaceholder, count)
	}

	return fmt.Sprintf("%s (%s)", c.Message, count)
}
