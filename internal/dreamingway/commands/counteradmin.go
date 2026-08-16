package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/kn-lim/chattingway/v2/counter"
	"github.com/kn-lim/dreamingway-bot/internal/constants"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
)

// counterListLimit is the largest reply that the list subcommand builds.
// Discord refuses a message of more than 2000 characters.
const counterListLimit = 1800

// counterAdmin handles the subcommands of /counter.
func counterAdmin(i discord.Interaction) (string, error) {
	guildID := i.GuildID()
	if guildID == nil {
		return "", errors.New("command must be sent from a server")
	}

	data := i.(discord.ApplicationCommandInteraction).SlashCommandInteractionData()
	if data.SubCommandName == nil {
		return "", errors.New("missing counter subcommand")
	}

	table := os.Getenv("COUNTER_TABLE_NAME")
	region := os.Getenv("AWS_REGION")

	// list only reads, so it needs no role. set and message change data.
	if *data.SubCommandName != "list" {
		allowed, err := hasCounterAdminRole(i)
		if err != nil {
			return "", err
		}
		if !allowed {
			return constants.Unauthorized, nil
		}
	}

	switch *data.SubCommandName {
	case "set": // /counter set <name> <value>
		c, err := counter.Set(context.TODO(), table, region, guildID.String(),
			data.String("name"), int64(data.Int("value")))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Set `%s` to `%d`.", c.Name, c.Count), nil
	case "message": // /counter message <name> <message>
		c, err := counter.SetMessage(context.TODO(), table, region, guildID.String(),
			data.String("name"), data.String("message"))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Set the message for `%s`. It now replies: %s",
			c.Name, renderCounter(c)), nil
	case "list": // /counter list
		counters, err := counter.List(context.TODO(), table, region, guildID.String())
		if err != nil {
			return "", err
		}

		return formatCounterList(counters), nil
	}

	return "", errors.New("invalid counter subcommand")
}

// hasCounterAdminRole reports whether the member holds the counter admin role.
func hasCounterAdminRole(i discord.Interaction) (bool, error) {
	member := i.Member()
	if member == nil {
		return false, errors.New("command must be sent from a server")
	}

	roles, err := dreamingway.FetchGuildRoles(i.GuildID(),
		os.Getenv("DISCORD_API_VERSION"), os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Name == os.Getenv("COUNTER_DISCORD_ADMIN_ROLE") {
			return slices.Contains(member.RoleIDs, role.ID), nil
		}
	}

	return false, nil
}

// formatCounterList builds the reply for the list subcommand.
// If the list is too long for Discord, formatCounterList stops and reports the rest.
func formatCounterList(counters []counter.Counter) string {
	if len(counters) == 0 {
		return "This server has no counters yet."
	}

	var b strings.Builder
	b.WriteString("**Counters**\n")

	for index, c := range counters {
		line := fmt.Sprintf("- `%s`: %d\n", c.Name, c.Count)
		if b.Len()+len(line) > counterListLimit {
			fmt.Fprintf(&b, "...and %d more.", len(counters)-index)
			break
		}

		b.WriteString(line)
	}

	return b.String()
}
