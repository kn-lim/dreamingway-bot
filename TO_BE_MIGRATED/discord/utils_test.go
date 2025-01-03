package discord_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kn-lim/dreamingway-bot/TO_BE_MIGRATED/discord"
)

const (
	TestMessage = "test message"
)

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestCheckRole(t *testing.T) {
	roles := []string{"1234567890", "test", "0987654321"}

	t.Run("success", func(t *testing.T) {
		if !discord.CheckRole(roles, "test") {
			t.Fatalf("discord.CheckRole() got %v, want %v", false, true)
		}
	})

	t.Run("failure", func(t *testing.T) {
		if discord.CheckRole(roles, "nonexistant") {
			t.Fatalf("discord.CheckRole() got %v, want %v", true, false)
		}
	})
}

func TestDeferredMessage(t *testing.T) {
	r, err := discord.DeferredMessage()

	if err != nil {
		t.Fatalf("discord.DeferredMessage() err = %v, want nil", err)
	}

	if r.Type != discordgo.InteractionResponseDeferredChannelMessageWithSource {
		t.Fatalf("discord.DeferredMessage() type = %v, want %v", r.Type, discordgo.InteractionResponseDeferredChannelMessageWithSource)
	}
}

func TestSendDeferredMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		if err := discord.SendDeferredMessage("appID", "token", TestMessage, discord.WithURL(server.URL)); err != nil {
			t.Fatalf("discord.SendDeferredMessage() err = %v, want nil", err)
		}
	})

	t.Run("non-200 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			_, err := w.Write([]byte(`{"error": "internal server error"}`))
			if err != nil {
				t.Fatalf("error! couldn't write to http")
			}
		}))
		defer server.Close()

		if err := discord.SendDeferredMessage("appID", "token", TestMessage, discord.WithURL(server.URL)); err == nil {
			t.Fatalf("discord.SendDeferredMessage() err = nil, want error")
		}
	})

	t.Run("non-200 response with error in body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		if err := discord.SendDeferredMessage("appID", "token", TestMessage, discord.WithURL(server.URL)); err == nil {
			t.Fatalf("discord.SendDeferredMessage() err = nil, want error")
		}
	})

	t.Run("error with http client", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		client := server.Client()

		// Create a mock RoundTripper to be used as Transport
		client.Transport = &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("error with mock transport")
			},
		}

		if err := discord.SendDeferredMessage("appID", "token", TestMessage, discord.WithClient(client), discord.WithURL(server.URL)); err == nil {
			t.Fatalf("discord.SendDeferredMessage() err = nil, want error")
		}
	})
}
