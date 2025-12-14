<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/dreamingway-bot/main/images/dreamingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# dreamingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/dreamingway-bot)
![GitHub Workflow Status - Build](https://img.shields.io/github/actions/workflow/status/kn-lim/dreamingway-bot/build.yaml)
![GitHub Workflow Status - Tests](https://img.shields.io/github/actions/workflow/status/kn-lim/dreamingway-bot/test.yaml?label=tests)
[![Coverage Status](https://coveralls.io/repos/github/kn-lim/dreamingway-bot/badge.svg?branch=main)](https://coveralls.io/github/kn-lim/dreamingway-bot?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kn-lim/dreamingway-bot)](https://goreportcard.com/report/github.com/kn-lim/dreamingway-bot)
![License](https://img.shields.io/github/license/kn-lim/dreamingway-bot)

A personal Discord bot to handle miscellaneous tasks hosted on AWS Lambda.

## Packages Used

- [aws-lambda-go](https://github.com/aws/aws-lambda-go/)
- [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)
- [chattingway](https://github.com/kn-lim/chattingway)
- [disgo](https://github.com/disgoorg/disgo)
- [koanf](https://github.com/knadh/koanf)
- [mergo](https://github.com/darccio/mergo)
- [urfave/cli](https://github.com/urfave/cli)
- [zap](https://github.com/uber-go/zap)

# Using the Discord Bot

## Discord Slash Commands

| Command | Description |
| - | - |
| `/coinflip` | Flips a coin |
| `/ping` | Ping |
| `/roll` | Rolls a dice with modifiers |

## How to Build

From the project home directory:

- **Endpoint Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`
- **Task Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/task/`

Zip the bootstrap binaries and upload it to the Lambda functions.

## Syncing Commands with Discord

1. Rename `config.example.json` to `config.json` and add in the values.
2. From the project directory, run the following command: `go run .`


```
NAME:
   dreamingway - sync discord commands

USAGE:
   dreamingway [global options]

GLOBAL OPTIONS:
   --verbose, -v               enable verbose logging (default: false)
   --config string, -c string  path to the configuration file (default: "config.json")
   --config-string string      configuration as a json string
   --help, -h                  show help
```

## Environment Variables

### Endpoint Lambda Function

| Name | Description |
| - | - |
| `DEBUG` | Enable debug mode |
| `TASK_FUNCTION_NAME` | Name of the Task Lambda Function |
| `DISCORD_BOT_APPLICATION_ID` | Discord Bot Application ID |
| `DISCORD_BOT_PUBLIC_KEY` | Discord Bot Public Key |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |

### Task Lambda Function

| Name | Description |
| - | - |
| `DISCORD_API_VERSION` | Discord API Version |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |
| `PZ_HOST` | Project Zomboid Host IP/URL |
| `PZ_HOST_INSTANCE_ID` | AWS Instance ID of the Project Zomboid Host |
| `PZ_HOST_REGION` | AWS Instance ID of the Project Zomboid Host |
| `PZ_RCON_PASSWORD` | RCON Password of the Project Zomboid server |

## AWS Setup

To quickly spin up **dreamingway-bot** on AWS, use the [Terraform module](https://github.com/kn-lim/chattingway-terraform/).

1. Create the **endpoint** Lambda function on AWS.
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `x86_64`.
2. Add an API Gateway triger to the **endpoint** Lambda function.
    - Use the following settings:
      - **Intent**: Create a new API
      - **API type**: REST API
      - **Security**: Open
3. Create the **task** Lambda function on AWS.
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `x86_64`.
4. Build the **endpoint** and **task** binaries.
5. Archive the `bootstrap` binaries in .zip files and upload it to the Lambda functions.
6. In the `Configuration` tab, add in the required environment variables to the Lambda functions.
7. Change the `Timeout` of the **task** Lambda function to a value greater than 3 seconds.
    - The `Timeout` of the **endpoint** Lambda function can stay as 3 seconds to follow Discord's requirements.

## Discord Setup

### Interactions Endpoint URL

Get the **endpoint** Lambda API Gateway triggers' `API endpoint` and add it to the Discord bot's `Interactions Endpoint URL` in the [Discord Developer Portal](https://discord.com/developers/).

### OAuth2 Scopes

In the OAuth2 URL Generator, give the following scopes when adding the bot to a server:

#### Scopes

- `applications.commands`
- `bot`

#### Bot Permissions

- `Manage Roles`
