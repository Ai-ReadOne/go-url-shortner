package main

import (
	"fmt"
	"os"

	"github.com/ai-readone/go-url-shortner/configs"
	"github.com/ai-readone/go-url-shortner/internal/app"
	"github.com/ai-readone/go-url-shortner/internal/database"
	"github.com/ai-readone/go-url-shortner/logger"
)

func main() {
	// initializes the configuration for the app,
	// using the config.yaml passed as argument to the go run command.
	config := configs.Init(os.Args[2])

	// setup connection with postgres server
	if err := database.InitPostgres(config); err != nil {
		panic(err)
	}

	router := app.RegisterRoutes(config)

	// starts the app's server
	logger.Info(fmt.Sprintf("Starting server %s", config.Server))
	logger.Fatal(router.Run(config.Server))
}
