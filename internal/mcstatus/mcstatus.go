package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MCStatusResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online int `json:"online"`
	} `json:"players"`
}

// GetMCStatus checks with mcstatus.io to get information about the Minecraft server
func GetMCStatus() (bool, int, error) {
	log.Println("GetMCStatus()")

	serverURL := fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
	if serverURL == "" {
		return false, 0, fmt.Errorf("PIXELMON_DOMAIN/PIXELMON_SUBDOMAIN environment variable not set")
	}

	url := fmt.Sprintf("https://api.mcstatus.io/v2/status/java/%s", serverURL)
	log.Printf("%v", url)

	response, err := http.Get(url)
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

	log.Printf("%v | Online: %v, Player Count: %v", serverURL, status.Online, status.Players.Online)

	return status.Online, status.Players.Online, nil
}
