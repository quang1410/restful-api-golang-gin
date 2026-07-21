package main

import (
	"galvin/lession05-exercise-user-management/internal/app"
	"galvin/lession05-exercise-user-management/internal/config"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
