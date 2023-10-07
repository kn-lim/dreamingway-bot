package pixelmon

import (
	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

func GetStatus() (bool, int, error) {
	// log.Println("GetStatus()")

	return mcstatus.GetMCStatus(ServerURL)
}

func StartInstance() error {
	return nil
}

func StartService() error {
	return nil
}

func StopInstance() error {
	return nil
}

func StopService() error {
	return nil
}
