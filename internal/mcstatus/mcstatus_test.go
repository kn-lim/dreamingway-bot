package mcstatus_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

func TestGetMCStatus(t *testing.T) {
	serverURL := "test-minecraft-server"

	t.Run("success", func(t *testing.T) {
		response := mcstatus.MCStatusResponse{
			Online: true,
			Players: struct {
				Online int `json:"online"`
			}{Online: 5},
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(response)
		}))
		defer server.Close()

		gotOnline, gotPlayers, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithBaseURL(server.URL))
		if err != nil {
			t.Fatalf("mcstatus.GetMCStatus() err = %v; want nil", err)
		}

		if !gotOnline {
			t.Fatalf("mcstatus.GetMCStatus() online = %v; want true", gotOnline)
		}

		if gotPlayers != 5 {
			t.Fatalf("mcstatus.GetMCStatus() players = %v; want %v", err, response.Players.Online)
		}
	})

	t.Run("error with http", func(t *testing.T) {
		invalidURL := "jttp://invalid-url"

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithBaseURL(invalidURL))
		if err == nil {
			t.Fatal("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})
}
