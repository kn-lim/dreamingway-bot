<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/dreamingway-bot/main/images/dreamingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# dreamingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/dreamingway-bot)
![GitHub Workflow Status - Build](https://img.shields.io/github/actions/workflow/status/kn-lim/dreamingway-bot/build.yaml)
![GitHub Workflow Status - Tests](https://img.shields.io/github/actions/workflow/status/kn-lim/dreamingway-bot/test.yaml?label=tests)
[![codecov](https://codecov.io/gh/kn-lim/dreamingway-bot/branch/main/graph/badge.svg)](https://codecov.io/gh/kn-lim/dreamingway-bot)
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

# How to Use the Discord Bot

## Discord Slash Commands

| Command | Description |
| - | - |
| `/coinflip` | Flips a coin |
| `/ping` | Ping command |
| `/roll` | Rolls dice with modifiers |

## Build

Run these commands from the project home directory:

- **Endpoint Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`
- **Task Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/task/`

Then compress each `bootstrap` binary into a separate ZIP file and upload each ZIP file to the related Lambda function.

## Sync Commands with Discord

1. Rename `config.example.json` to `config.json`.
2. Add the values to `config.json`.
3. Run this command from the project directory: `go run .`

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
| `DISCORD_BOT_APPLICATION_ID` | Discord Bot Application ID |
| `DISCORD_BOT_PUBLIC_KEY` | Discord Bot Public Key |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |
| `TASK_FUNCTION_NAME` | Name of the Task Lambda Function |

### Task Lambda Function

| Name | Description |
| - | - |
| `CLOUDFLARE_API_TOKEN` | Cloudflare API Token |
| `CLOUDFLARE_ZONE_ID` | Cloudflare Zone ID |
| `DISCORD_API_VERSION` | Discord API Version |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |
| `PZ_DISCORD_ADMIN_ROLE` | Discord Admin Role for the Project Zomboid |
| `PZ_HOST` | Project Zomboid Host IP/URL |
| `PZ_HOST_INSTANCE_ID` | AWS Instance ID of the Project Zomboid Host |
| `PZ_HOST_REGION` | AWS Region of the Project Zomboid Host |
| `PZ_RCON_PASSWORD` | RCON Password of the Project Zomboid server |
| `PZ_RCON_PORT` | RCON Port of the Project Zomboid server |

## AWS Setup

### Terraform

To create **dreamingway-bot** on AWS, use the [Terraform module](https://github.com/kn-lim/terraform-aws-chattingway/).

### Manual

1. Create the **endpoint** Lambda function on AWS.
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `arm64`.
2. Add an API Gateway trigger to the **endpoint** Lambda function.
    - Use the following settings:
      - **Intent**: Create a new API
      - **API type**: HTTP API
      - **Security**: Open
3. Create the **task** Lambda function on AWS.
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `arm64`.
4. Build the **endpoint** and **task** binaries. Name each binary `bootstrap`.
5. Compress each `bootstrap` binary into a separate `bootstrap.zip` file.
6. Upload each ZIP file to the related Lambda function.
7. In the `Configuration` tab, add the required environment variables to each Lambda function.
8. Change the `Timeout` of the **task** Lambda function to more than 3 seconds.
    - Keep the `Timeout` of the **endpoint** Lambda function at 3 seconds. Discord requires a response in 3 seconds.

## Discord Setup

### Interactions Endpoint URL

1. Find the `API endpoint` value of the API Gateway trigger on the **endpoint** Lambda function.
2. Add this value to the `Interactions Endpoint URL` field of the bot in the [Discord Developer Portal](https://discord.com/developers/).

### OAuth2 Scopes

When you add the bot to a server, use the OAuth2 URL Generator. Give the bot these scopes and permissions:

#### Scopes

- `applications.commands`
- `bot`

#### Bot Permissions

- `Manage Roles`
