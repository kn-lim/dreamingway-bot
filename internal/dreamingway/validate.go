package dreamingway

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

// ValidateRequest validates a Discord interaction request
func ValidateRequest(request events.APIGatewayProxyRequest, publicKey string) error {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}

	var body []byte
	if request.IsBase64Encoded {
		body_bytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return err
		}

		body = body_bytes
	} else {
		body = []byte(request.Body)
	}

	ed25519PublicKey := ed25519.PublicKey(publicKeyBytes)

	x_signature, ok := request.Headers["x-signature-ed25519"]
	if !ok {
		return errors.New("missing x-signature-ed25519 header")
	}

	x_signature_time, ok := request.Headers["x-signature-timestamp"]
	if !ok {
		return errors.New("missing x-signature-timestamp header")
	}

	x_signature_bytes, err := hex.DecodeString(x_signature)
	if err != nil {
		return err
	}

	signed_data := []byte(x_signature_time + string(body))

	if !ed25519.Verify(ed25519PublicKey, signed_data, x_signature_bytes) {
		return errors.New("invalid request signature")
	}

	return nil
}
