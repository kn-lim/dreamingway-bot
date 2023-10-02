package pixelmon

import (
	"log"

	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

func GetStatus() (bool, int, error) {
	log.Println("GetStatus()")

	return mcstatus.GetMCStatus()
}
