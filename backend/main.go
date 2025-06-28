package main

import (
	"log"

	"reminder-app/app"
)

func main() {
	application := app.New()
	if err := application.Run(); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}