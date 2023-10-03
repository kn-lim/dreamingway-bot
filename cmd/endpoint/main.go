package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/discord"
)

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	public_key_bytes, err := hex.DecodeString(os.Getenv("DISCORD_BOT_PUBLIC_KEY"))
	if err != nil {
		return events.LambdaFunctionURLResponse{}, errors.New("error! couldn't decode public key")
	}
	if request.Body == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"error": "Body is empty"}`,
		}, errors.New("error! body is empty")
	}

	var body []byte

	if request.IsBase64Encoded {
		body_bytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.LambdaFunctionURLResponse{}, fmt.Errorf("error! couldn't decode body [%s]: %s", body, err)
		}

		body = body_bytes
	} else {
		body = []byte(request.Body)
	}

	public_key := ed25519.PublicKey(public_key_bytes)

	x_signature, ok := request.Headers["x-signature-ed25519"]
	if !ok {
		log.Print("Received Signature Header Error (400)")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"error": "Missing x-signature-ed25519 header"}`,
		}, errors.New("error! missing x-signature-ed25519 header")
	}

	x_signature_time, ok := request.Headers["x-signature-timestamp"]
	if !ok {
		log.Print("Received Timestamp Header Error (400)")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"error": "Missing x-signature-timestamp header"}`,
		}, errors.New("error! missing x-signature-timestamp header")
	}

	x_signature_bytes, err := hex.DecodeString(x_signature)
	if err != nil {
		return events.LambdaFunctionURLResponse{}, errors.New("error! couldn't decode signature")
	}

	signed_data := []byte(x_signature_time + string(body))

	if !ed25519.Verify(public_key, signed_data, x_signature_bytes) {
		// Unauthorized access
		// log.Print("Received Unauthorized (401)")
		return events.LambdaFunctionURLResponse{
			StatusCode: 401,
		}, nil
	} else {
		// Authorized access
		var interaction discordgo.Interaction
		if err := json.Unmarshal(body, &interaction); err != nil {
			log.Printf("Error! Could not decode interaction: %s", err)
			return events.LambdaFunctionURLResponse{
				StatusCode: 400,
			}, nil
		}

		switch {
		case interaction.Type == 1:
			// Ping (200)
			log.Print("Received Ping (200)")
			return events.LambdaFunctionURLResponse{
				StatusCode: 200,
				Body:       `{"type": 1}`,
			}, nil
		case interaction.Type == 2:
			// Application Command
			log.Printf("Recieved Application Command: %s", interaction.ApplicationCommandData().Name)

			cmd, ok := discord.Commands[interaction.ApplicationCommandData().Name]
			if !ok {
				// log.Printf("Error! Command does not exist: %s", interaction.ApplicationCommandData().Name)
				return events.LambdaFunctionURLResponse{
					StatusCode: 404,
					Body:       `{"error": "Command does not exist"}`,
				}, nil
			}

			var response discordgo.InteractionResponse
			if cmd.Handler != nil {
				response, err = cmd.Handler(&interaction)
			} else if cmd.Options[interaction.ApplicationCommandData().Options[0].Name] != nil {
				response, err = cmd.Options[interaction.ApplicationCommandData().Options[0].Name](&interaction)
			}

			if err != nil {
				log.Printf("Error! Handler Error: %s", err)
				return events.LambdaFunctionURLResponse{}, err
			} else {
				response_body, err := json.Marshal(&response)
				if err != nil {
					log.Printf("Error! Couldn't marshal JSON: %s", err)
					return events.LambdaFunctionURLResponse{}, err
				}
				// log.Printf("Response Body: %s", string(response_body))

				return events.LambdaFunctionURLResponse{
					StatusCode: 200,
					Body:       string(response_body),
				}, nil
			}
		default:
			// Unknown (501)
			// log.Print("Received Unknown (501)")
			return events.LambdaFunctionURLResponse{
				StatusCode: 501,
			}, nil
		}
	}
}

func main() {
	lambda.Start(handler)
}
