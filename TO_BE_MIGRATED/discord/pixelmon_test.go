package discord

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestStatus(t *testing.T) {
	t.Run("success with online server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{"online": true, "players": {"online": 5}}`))
			if err != nil {
				t.Fatalf("error! couldn't write to http")
			}
		}))
		defer server.Close()

		message, err := status(&discordgo.Interaction{}, WithURL(server.URL))

		if err != nil {
			t.Fatalf("discord.status() err = %v, want nil", err)
		}

		if message == "" {
			t.Fatalf("discord.status() message is nil, want non-nil")
		}
	})

	t.Run("success with offline server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{"online": false, "players": {"online": 0}}`))
			if err != nil {
				t.Fatalf("error! couldn't write to http")
			}
		}))
		defer server.Close()

		message, err := status(&discordgo.Interaction{}, WithURL(server.URL))

		if err != nil {
			t.Fatalf("discord.status() err = %v, want nil", err)
		}

		if message == "" {
			t.Fatalf("discord.status() message is nil, want non-nil")
		}
	})

	t.Run("error with pixelmon.GetStatus()", func(t *testing.T) {
		invalidURL := "http://invalid-url"

		message, err := status(&discordgo.Interaction{}, WithURL(invalidURL))

		if err == nil {
			t.Fatalf("discord.status() err = nil, want non-nil")
		}

		if message != "" {
			t.Fatalf("discord.status() message = %v, want nil", message)
		}
	})
}
