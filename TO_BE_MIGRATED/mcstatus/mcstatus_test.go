package mcstatus_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kn-lim/dreamingway-bot/TO_BE_MIGRATED/mcstatus"
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

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				t.Fatalf("error! couldn't encode json")
			}
		}))
		defer server.Close()

		gotOnline, gotPlayers, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithURL(server.URL))
		if err != nil {
			t.Fatalf("mcstatus.GetMCStatus() err = %v; want nil", err)
		}

		if !gotOnline {
			t.Fatalf("mcstatus.GetMCStatus() online = %v; want %v", gotOnline, true)
		}

		if gotPlayers != 5 {
			t.Fatalf("mcstatus.GetMCStatus() players = %v; want %v", err, response.Players.Online)
		}
	})

	t.Run("error with http", func(t *testing.T) {
		invalidURL := "http://invalid-url"

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithURL(invalidURL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})

	t.Run("error with reading body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Use a hijacker to forcibly close the connection
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("error! couldn't hijack connection")
			}
			conn, _, err := hijacker.Hijack()
			if err != nil {
				t.Fatal(err)
			}
			conn.Close() // Close connection to simulate read error
		}))
		defer server.Close()

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithURL(server.URL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})

	t.Run("error with json unmarshalling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Send an invalid JSON response
			_, err := w.Write([]byte("invalid json"))
			if err != nil {
				t.Fatalf("error! couldn't write to http")
			}
		}))
		defer server.Close()

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithURL(server.URL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})
}
