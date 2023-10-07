package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type MCStatusResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online int `json:"online"`
	} `json:"players"`
}

// GetMCStatus checks with mcstatus.io to get information about the Minecraft server
func GetMCStatus(url string) (bool, int, error) {
	// log.Println("GetMCStatus()")

	mcstatus := fmt.Sprintf("https://api.mcstatus.io/v2/status/java/%s", url)
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

	log.Printf("%v | Online: %v, Player Count: %v", url, status.Online, status.Players.Online)

	return status.Online, status.Players.Online, nil
}
