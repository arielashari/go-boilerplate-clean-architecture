package main

import (
	"log"
	"os"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/app"
	"github.com/spf13/viper"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	viper.Set("APP_ENV", env)

	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application := app.NewApp(cfg)

	defer application.DatabaseInfra.CloseAll()

	application.Server.Start()
}
