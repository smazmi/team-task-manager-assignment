package main

import (
	"log"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application startup failed: %v", err)
	}
}
