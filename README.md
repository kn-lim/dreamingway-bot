<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/dreamingway-bot/main/images/dreamingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# dreamingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/dreamingway-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/kn-lim/dreamingway-bot)](https://goreportcard.com/report/github.com/kn-lim/dreamingway-bot)
![License](https://img.shields.io/github/license/kn-lim/dreamingway-bot)

A personal Discord bot to handle miscellaneous tasks hosted on AWS Lambda.

## Packages Used

- [aws-lambda-go](https://github.com/aws/aws-lambda-go/)
- [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)
- [discordgo](https://github.com/bwmarrin/discordgo/)

# Using the Discord Bot

## How to Build

From the project home directory: 

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`

## Syncing Commands with Discord

https://github.com/kn-lim/dreamingway-bot/tree/main/cmd/cli

## Environment Variables

### Endpoint Lambda Function

#### AWS

| Name | Description |
| - | - |
| `TASK_FUNCTION_NAME` | Name of the Task Lambda Function |

#### Discord

| Name | Description |
| - | - |
| `DISCORD_BOT_APPLICATION_ID` | Discord Bot Application ID |
| `DISCORD_BOT_PUBLIC_KEY` | Discord Bot Public Key |
| `DISCORD_BOT_TOKEN` | Discord Bot Token |

### Task Lambda Function

#### Discord

| Name | Description |
| - | - |
| `DISCORD_API_VERSION` | Discord API Version |

#### Commands

| Name | Description |
| - | - |
| `PIXELMON_ROLE_ID` | Role ID to allow `/pixelmon` command |

#### Pixelmon

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

1. Create a Lambda endpoint function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
    - Under `Advanced Settings`, select:
        - `Enable function URL`
          - Auth type: `NONE`
          - Invoke mode: `BUFFERED (default)`
          - Enable `Configure cross-origin resource sharing (CORS)`
2. Create a Lambda task function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
    - Under `Advanced Settings`, select:
        - `Enable VPC`
3. Archive the `bootstrap` binary in a .zip file and upload it to the Lambda functions.
4. In the `Configuration` tab, add in the required environment variables.
5. Get the Lambda endpoint function's `Function URL` and add it to the Discord bot's `Interactions Endpoint URL` in the [Discord Developer Portal](https://discord.com/developers/).
    - If it saves properly, that indicates your Lambda function is properly configured to act as a Discord bot.
