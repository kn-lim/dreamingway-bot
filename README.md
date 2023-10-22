<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/dreamingway-bot/main/images/dreamingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# dreamingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/dreamingway-bot)
[![Coverage Status](https://coveralls.io/repos/github/kn-lim/dreamingway-bot/badge.svg?branch=main)](https://coveralls.io/github/kn-lim/dreamingway-bot?branch=main)
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

- **Endpoint Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`
- **Task Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/task/`

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
| `DISCORD_BOT_TOKEN` | Discord Bot Token |

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
| `PIXELMON_ROLE_ID` | Role ID to allow `/pixelmon` command |

## AWS Setup

1. Create an **endpoint** Lambda function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
    - Under `Advanced Settings`, select:
        - `Enable function URL`
          - Auth type: `NONE`
          - Invoke mode: `BUFFERED (default)`
          - Enable `Configure cross-origin resource sharing (CORS)`
2. Create a **task** Lambda function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
3. Archive the `bootstrap` binary in a .zip file and upload it to the Lambda functions.
4. In the `Configuration` tab, add in the required environment variables to the Lambda functions.
5. Give the role the **task** Lambda function is using permission to access the AWS resources it will need.
6. Change the `Timeout` of the **task** Lambda function to a value greater than 3 seconds.
    - The `Timeout` of the **endpoint** Lambda function can stay as 3 seconds to follow Discord's requirements.
8. Get the **endpoint** Lambda function's `Function URL` and add it to the Discord bot's `Interactions Endpoint URL` in the [Discord Developer Portal](https://discord.com/developers/).
    - If it saves properly, that indicates your Lambda function is properly configured to act as a Discord bot.

