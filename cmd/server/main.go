package main

import (
	"log"

	"go_web/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap application failed: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("run application failed: %v", err)
	}
}
