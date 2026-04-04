package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

// -ldflags="-X main.VERSION=X.X.X"
var VERSION = "0.0.0"

func main() {
	godotenv.Load()

	config, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load configuration: %s\n", err)
		os.Exit(1)
	}

	/***********************/

	logOptions := &slog.HandlerOptions{
		Level: config.LogLevel,
	}

	var log *slog.Logger
	if VERSION == "0.0.0" {
		log = slog.New(slog.NewTextHandler(os.Stdout, logOptions))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, logOptions))
	}

	slog.SetDefault(log)

	log.Info("Hello World", "version", VERSION)

	/***********************/

	os.Exit(0)
}
