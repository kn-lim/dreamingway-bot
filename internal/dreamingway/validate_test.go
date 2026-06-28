package dreamingway

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestValidateRequest_ValidSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKey := hex.EncodeToString(pub)

	timestamp := "1234567890"
	body := `{"type":1}`
	signature := ed25519.Sign(priv, []byte(timestamp+body))

	req := events.APIGatewayV2HTTPRequest{
		Body: body,
		Headers: map[string]string{
			"x-signature-ed25519":   hex.EncodeToString(signature),
			"x-signature-timestamp": timestamp,
		},
	}

	if err := ValidateRequest(req, publicKey); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestValidateRequest_InvalidSignature(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKey := hex.EncodeToString(pub)

	req := events.APIGatewayV2HTTPRequest{
		Body: `{"type":1}`,
		Headers: map[string]string{
			"x-signature-ed25519":   hex.EncodeToString(make([]byte, ed25519.SignatureSize)),
			"x-signature-timestamp": "1234567890",
		},
	}

	if err := ValidateRequest(req, publicKey); err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}
}

func TestValidateRequest_MissingSignatureHeader(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKey := hex.EncodeToString(pub)

	req := events.APIGatewayV2HTTPRequest{
		Body:    `{"type":1}`,
		Headers: map[string]string{"x-signature-timestamp": "1234567890"},
	}

	if err := ValidateRequest(req, publicKey); err == nil {
		t.Fatal("expected error for missing signature header, got nil")
	}
}
