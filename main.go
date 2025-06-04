package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
			// Create a new Discord session
			d, err := dreamingway.NewDreamingway(cfg.Token)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to create discord session: %v", err), 1)
			}

			applicationID, err := snowflake.Parse(cfg.AppID)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to parse application ID: %v", err), 1)
			}

			// Sync global commands
			if err := commands.SyncGlobalCommands(d.Client.Rest(), applicationID, commands.GlobalCommands); err != nil {
				return cli.Exit(fmt.Sprintf("failed to sync global commands: %v", err), 1)
			}

			if !cmd.Bool("verbose") {
				fmt.Println("Global commands synced successfully.")
			}

			for _, server := range cfg.Guilds {
				snowflakeID, err := snowflake.Parse(server.GuildID)
				if err != nil {
					return cli.Exit(fmt.Sprintf("failed to parse guild ID: %v", err), 1)
				}
				// Check for server specific commands
				if len(server.Commands) == 0 {
					continue
				}

				// Sync guild commands
				cmds := make(map[string]commands.Command)
				for key, serverCmd := range commands.Commands {
					for _, cmd := range server.Commands {
						if serverCmd.Command.CommandName() == cmd {
							cmds[key] = serverCmd
						}
					}
				}
				if err := commands.SyncGuildCommands(d.Client.Rest(), applicationID, snowflakeID, cmds); err != nil {
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
