package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type options struct {
	baseURL string
}
type Option func(*options)

func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		o.baseURL = baseURL
	}
}

type MCStatusResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online int `json:"online"`
	} `json:"players"`
}

// GetMCStatus checks with mcstatus.io to get information about the Minecraft server
func GetMCStatus(serverURL string, opts ...Option) (bool, int, error) {
	// log.Println("GetMCStatus()")

	// Defaults
	config := &options{
		baseURL: "https://api.mcstatus.io/v2/status/java",
	}
	for _, opt := range opts {
		opt(config)
	}

	mcstatus := fmt.Sprintf("%s/%s", config.baseURL, serverURL)
	// log.Printf("MCStatus URL: %v", mcstatus)

	response, err := http.Get(mcstatus)
	if err != nil {
		return false, 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, 0, err
	}

	var status MCStatusResponse
	err = json.Unmarshal(body, &status)
	if err != nil {
		return false, 0, err
	}

	// log.Printf("%v | Online: %v, Player Count: %v", url, status.Online, status.Players.Online)

	return status.Online, status.Players.Online, nil
}
