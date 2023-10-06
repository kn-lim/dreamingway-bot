package discord

import "log"

func ping() (string, error) {
	log.Println("/ping")

	return "Pong!", nil
}
