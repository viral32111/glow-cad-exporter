package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

const (
	VERSION = "0.1.0"
)

func main() {
	godotenv.Load()

	config, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load configuration: %s\n", err)
		os.Exit(1)
	}

	/***********************/

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	}))
	slog.SetDefault(log)

	log.Info("Hello", "version", VERSION)

	/***********************/

	os.Exit(0)
}
