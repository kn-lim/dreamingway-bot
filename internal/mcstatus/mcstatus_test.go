package mcstatus_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

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
		invalidURL := "http://invalid-url"

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithBaseURL(invalidURL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})

	t.Run("error with reading body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Using a hijacker to forcibly close the connection
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("error! couldn't hijack connection")
			}
			conn, _, err := hijacker.Hijack()
			if err != nil {
				t.Fatal(err)
			}
			conn.Close() // Forcibly close connection to simulate read error
		}))
		defer server.Close()

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithBaseURL(server.URL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})

	t.Run("error with json unmarshalling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Sending an invalid JSON response
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		_, _, err := mcstatus.GetMCStatus(serverURL, mcstatus.WithBaseURL(server.URL))
		if err == nil {
			t.Fatalf("mcstatus.GetMCStatus() err = nil; want non-nil")
		}
	})
}
