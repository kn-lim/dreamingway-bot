package dreamingway

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

// ValidateRequest validates a Discord interaction request
func ValidateRequest(request events.APIGatewayV2HTTPRequest, publicKey string) error {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}

	var body []byte
	if request.IsBase64Encoded {
		bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return err
		}

		body = bodyBytes
	} else {
		body = []byte(request.Body)
	}

	ed25519PublicKey := ed25519.PublicKey(publicKeyBytes)

	signature, ok := request.Headers["x-signature-ed25519"]
	if !ok {
		return errors.New("missing x-signature-ed25519 header")
	}

	timestamp, ok := request.Headers["x-signature-timestamp"]
	if !ok {
		return errors.New("missing x-signature-timestamp header")
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}

	signedData := []byte(timestamp + string(body))

	if !ed25519.Verify(ed25519PublicKey, signedData, signatureBytes) {
		return errors.New("invalid request signature")
	}

	return nil
}
