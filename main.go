package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/disgoorg/snowflake/v2"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v3"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/dreamingway/commands"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

var (
	k = koanf.New(".")

	cfg Config
)

type Config struct {
	AppID          string   `koanf:"app_id"`
	Token          string   `koanf:"token"`
	GlobalCommands []string `koanf:"global_commands"`
	Guilds         []struct {
		GuildID  string   `koanf:"guild_id"`
		Commands []string `koanf:"commands"`
	} `koanf:"guilds"`
}

// guildSync carries one guild's parsed ID and filtered commands from the validation phase to the sync phase.
type guildSync struct {
	guildID  snowflake.ID
	commands map[string]commands.Command
}

func main() {
	cmd := &cli.Command{
		Name:  "dreamingway",
		Usage: "sync discord commands",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "enable verbose logging",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to the configuration file",
				Value:   "config.json",
			},
			&cli.StringFlag{
				Name:  "config-string",
				Usage: "configuration as a json string",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "show the commands to publish, without any change to discord",
				Value: false,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Initialize logger
			var err error
			utils.Logger, err = utils.NewLogger(cmd.Bool("verbose"))
			if err != nil {
				return ctx, fmt.Errorf("couldn't initialize logger: %w", err)
			}

			// Check if config string is provided
			if cmd.String("config-string") != "" {
				// Parse the config string
				if err := k.Load(rawbytes.Provider([]byte(cmd.String("config-string"))), json.Parser()); err != nil {
					return ctx, fmt.Errorf("failed to load config string: %w", err)
				}
			} else {
				// Read the config file
				configFilePath := cmd.String("config")
				if _, err := os.Stat(configFilePath); err != nil {
					return ctx, fmt.Errorf("config file not found: %w", err)
				}
				if err := k.Load(file.Provider(configFilePath), json.Parser()); err != nil {
					return ctx, fmt.Errorf("failed to load config file: %w", err)
				}
			}

			if err := k.Unmarshal("", &cfg); err != nil {
				return ctx, fmt.Errorf("failed to unmarshal config: %w", err)
			}

			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			applicationID, err := snowflake.Parse(cfg.AppID)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to parse application ID: %v", err), 1)
			}

			// Validate and select global commands from the config.
			globalCmds, err := commands.FilterCommands(commands.GlobalCommands, cfg.GlobalCommands)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to select global commands: %v", err), 1)
			}

			// Validate and select every guild's commands from the config.
			// The CLI checks all guilds before it writes anything to Discord.
			guildSyncs := make([]guildSync, 0, len(cfg.Guilds))
			for _, server := range cfg.Guilds {
				snowflakeID, err := snowflake.Parse(server.GuildID)
				if err != nil {
					return cli.Exit(fmt.Sprintf("failed to parse guild ID: %v", err), 1)
				}

				guildCmds, err := commands.FilterCommands(commands.AllCommands(), server.Commands)
				if err != nil {
					return cli.Exit(fmt.Sprintf("failed to select commands for guild %s: %v", server.GuildID, err), 1)
				}

				guildSyncs = append(guildSyncs, guildSync{guildID: snowflakeID, commands: guildCmds})
			}

			// All config values are valid now. A dry run stops here and prints a report.
			if cmd.Bool("dry-run") {
				printDryRunReport(globalCmds, guildSyncs)
				return nil
			}

			// Create a new Discord session
			d, err := dreamingway.NewDreamingway(cfg.Token)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to create discord session: %v", err), 1)
			}

			if err := commands.SyncGlobalCommands(d.Client.Rest, applicationID, globalCmds); err != nil {
				return cli.Exit(fmt.Sprintf("failed to sync global commands: %v", err), 1)
			}

			if !cmd.Bool("verbose") {
				fmt.Println("Global commands synced successfully.")
			}

			for _, gs := range guildSyncs {
				if err := commands.SyncGuildCommands(d.Client.Rest, applicationID, gs.guildID, gs.commands); err != nil {
					return cli.Exit(fmt.Sprintf("failed to sync guild commands: %v", err), 1)
				}
			}

			if !cmd.Bool("verbose") {
				fmt.Println("Guild commands synced successfully.")
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// printDryRunReport prints alphabetically the commands a sync will publish and delete, without contacting Discord.
func printDryRunReport(globalCmds map[string]commands.Command, guildSyncs []guildSync) {
	fmt.Println("Dry run. Discord does not change.")
	fmt.Println()

	globalNames := sortedCommandNames(globalCmds)
	fmt.Printf("Global commands (%d):\n", len(globalNames))
	if len(globalNames) > 0 {
		for _, name := range globalNames {
			fmt.Printf("  /%s\n", name)
		}
	}

	if len(guildSyncs) == 0 {
		fmt.Println()
		fmt.Println("The config has no guilds.")
		return
	}

	for _, gs := range guildSyncs {
		guildNames := sortedCommandNames(gs.commands)

		fmt.Println()
		fmt.Printf("Guild %s (%d):\n", gs.guildID, len(guildNames))
		if len(guildNames) > 0 {
			for _, name := range guildNames {
				fmt.Printf("  /%s\n", name)
			}
		}
	}
}

// sortedCommandNames returns the map's keys in alphabetical order.
func sortedCommandNames(cmds map[string]commands.Command) []string {
	names := make([]string, 0, len(cmds))
	for name := range cmds {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}
