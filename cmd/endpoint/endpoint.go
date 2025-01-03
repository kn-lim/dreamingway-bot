package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/bwmarrin/discordgo"

	"github.com/kn-lim/dreamingway-bot/internal/dreamingway"
	"github.com/kn-lim/dreamingway-bot/internal/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize logger
	var err error
	utils.Logger, err = utils.NewLogger()
	if err != nil {
		log.Printf("couldn't initialize logger: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// Validate the request
	if err := dreamingway.ValidateRequest(request, os.Getenv("DISCORD_BOT_PUBLIC_KEY")); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
		}, err
	}

	// Parse the request body
	var body []byte
	if request.IsBase64Encoded {
		body_bytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		body = body_bytes
	} else {
		body = []byte(request.Body)
	}

	// Get discord interaction
	var interaction discordgo.Interaction
	if err := json.Unmarshal(body, &interaction); err != nil {
		utils.Logger.Errorw("couldn't unmarshal interaction",
			"error", err,
			"body", string(body),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// Handle the interaction
	switch interaction.Type {
	// Ping interaction
	case discordgo.InteractionPing:
		utils.Logger.Info("received ping interaction")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"type": 1}`,
		}, nil
	// Application command interaction
	case discordgo.InteractionApplicationCommand:
		// Get username of the user who sent the interaction
		username := "???"
		if interaction.Member != nil {
			username = interaction.Member.User.Username
		} else if interaction.User != nil {
			username = interaction.User.Username
		}

		utils.Logger.Infow("received application command interaction",
			"command", interaction.ApplicationCommandData().Name,
			"user", username,
		)

		// Get deferred response
		deferredResponse, err := json.Marshal(dreamingway.DeferredMessage())
		if err != nil {
			utils.Logger.Errorw("couldn't marshal deferred response",
				"error", err,
			)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		// Invoke task function
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
		if err != nil {
			utils.Logger.Errorw("couldn't load default config",
				"error", err,
			)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		client := lambdaSvc.NewFromConfig(cfg)
		payload, err := json.Marshal(&interaction)
		if err != nil {
			utils.Logger.Errorw("couldn't marshal interaction",
				"error", err,
			)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		input := &lambdaSvc.InvokeInput{
			FunctionName:   aws.String(os.Getenv("TASK_FUNCTION_NAME")),
			Payload:        payload,
			InvocationType: types.InvocationTypeEvent,
		}
		if _, err := client.Invoke(ctx, input); err != nil {
			utils.Logger.Errorw("couldn't invoke task function",
				"error", err,
			)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		// Return deferred response
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(deferredResponse),
		}, nil
	// Unknown interaction
	default:
		utils.Logger.Errorw("unsupported interaction type",
			"type", interaction.Type,
		)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
