<p align="center">
  <img width="100" src="https://raw.githubusercontent.com/kn-lim/dreamingway-bot/main/images/dreamingway.png"></img>
</p>

# dreamingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/dreamingway-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/kn-lim/dreamingway-bot)](https://goreportcard.com/report/github.com/kn-lim/dreamingway-bot)
![License](https://img.shields.io/github/license/kn-lim/dreamingway-bot)

A personal Discord bot to handle miscellaneous tasks hosted on AWS Lambda.

## Packages Used

- [aws-lambda-go](https://github.com/aws/aws-lambda-go/)
- [discordgo](https://github.com/bwmarrin/discordgo/)

# Using the Discord Bot

## How to Build

From the project home directory: 

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`

## Syncing Commands with Discord

https://github.com/kn-lim/dreamingway-bot/tree/main/cmd/cli

## Environment Variables

### Discord Specific

| Name | Description |
| - | - |
| `DISCORD_BOT_APPLICATION_ID` | Discord Bot Application ID |
| `DISCORD_BOT_PUBLIC_KEY` | Discord Bot Public Key |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |

### Commands

| Name | Description |
| - | - |
| `PIXELMON_ROLE_ID` | Role ID to allow `/pixelmon` command |

### Pixelmon Specific

| Name | Description |
| - | - |
| `PIXELMON_NAME` | AWS Name Tag of Pixelmon EC2 Instance |
| `PIXELMON_INSTANCE_ID` | AWS Instance ID of Pixelmon EC2 Instance |
| `PIXELMON_REGION` | AWS Region of Pixelmon EC2 Instance |
| `PIXELMON_HOSTED_ZONE_ID` | AWS Hosted Zone ID of Domain |
| `PIXELMON_DOMAIN` | Domain of Pixelmon Server |
| `PIXELMON_SUBDOMAIN` | Subdomain of Pixelmon Server |
| `PIXELMON_RCON_PASSWORD` | RCON Password of Pixelmon Service |

## AWS Setup

1. Create a Lambda function on AWS.
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
    - Under `Advanced Settings`, select:
        - `Enable function URL`
        - `Enable VPC` 
2. Archive the `bootstrap` binary in a .zip file and upload it to the Lambda function.
3. In the `Configuration` tab, add in the required environment variables.
4. Get the Lambda function's `Function URL` and add it to the Discord bot's `Interactions Endpoint URL` in the [Discord Developer Portal](https://discord.com/developers/).
    - If it saves properly, that indicates your Lambda function is properly configured to act as a Discord bot.
