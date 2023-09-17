package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	public_key_bytes, err := hex.DecodeString(os.Getenv("DISCORD_BOT_PUBLIC_KEY"))
	if err != nil {
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	}
	if request.Body == "" {
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	}

	var body []byte

	if request.IsBase64Encoded {
		body_bytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			// TODO
			return events.APIGatewayProxyResponse{}, errors.New("Error!")
		}

		body = body_bytes
	} else {
		body = []byte(request.Body)
	}

	public_key := ed25519.PublicKey(public_key_bytes)

	x_signature, ok := request.Headers["x-signature-ed25519"]
	if !ok {
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	}

	x_signature_time, ok := request.Headers["x-signature-timestamp"]
	if !ok {
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	}

	x_signature_bytes, err := hex.DecodeString(x_signature)
	if err != nil {
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	}

	signed_data := []byte(x_signature_time + string(body))

	if !ed25519.Verify(public_key, signed_data, x_signature_bytes) {
		// Unauthorized access
		// TODO
		return events.APIGatewayProxyResponse{}, errors.New("Error!")
	} else {
		// Authorized access
		// TODO

		var interaction discordgo.Interaction
		if err := json.Unmarshal(body, &interaction); err != nil {
			// TODO
			return events.APIGatewayProxyResponse{}, errors.New("Error!")
		}

		switch {
		case interaction.Type == 1:
			// Interaction: Ping (200)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"type": 1}`,
			}, nil
		case interaction.Type == 2:
			// Interaction: Application Command
			return events.APIGatewayProxyResponse{}, nil
		default:
			// TODO
			return events.APIGatewayProxyResponse{
				StatusCode: 501,
			}, nil
		}
	}
}

func main() {
	lambda.Start(handler)
}
