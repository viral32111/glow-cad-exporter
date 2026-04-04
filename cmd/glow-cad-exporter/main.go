package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/viral32111/glow-cad-exporter/internal/broker"
	"github.com/viral32111/glow-cad-exporter/internal/config"
	"github.com/viral32111/glow-cad-exporter/internal/metrics"
	"github.com/viral32111/glow-cad-exporter/internal/server"
)

// -ldflags="-X main.VERSION=X.X.X"
var VERSION = "0.0.0"

func main() {
	isLocalDevelopment := false
	if VERSION == "0.0.0" {
		isLocalDevelopment = true
	}

	/***********************/

	godotenv.Load()

	config, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load configuration: %s\n", err)
		os.Exit(1)
	}

	/***********************/

	var log *slog.Logger
	logOptions := &slog.HandlerOptions{Level: config.LogLevel}
	if isLocalDevelopment {
		log = slog.New(slog.NewTextHandler(os.Stdout, logOptions))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, logOptions))
	}

	slog.SetDefault(log)

	log.Info("Hello World", "version", VERSION)

	/***********************/

	/*
		signals := make(chan os.Signal, 1)
		isDone := make(chan bool, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-signals

			slog.Info("Caught termination signal, stopping...")
			isDone <- true
		}()
	*/

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	executionGroup, ctx := errgroup.WithContext(ctx)

	/***********************/

	httpServer, err := server.New(config)
	if err != nil {
		log.Error("Unable to create HTTP server", "error", err)
		os.Exit(1)
	}

	// Creating a new registry discards Go built-in metrics (runtime, process, etc.)
	prometheusRegistry := prometheus.NewRegistry()
	metrics.Register(prometheusRegistry)

	// Serve metrics on the configured path
	httpServer.HTTPMultiplexer.Handle(config.MetricsPath, promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{
		Registry: prometheusRegistry,

		EnableOpenMetrics: true,

		DisableCompression: false,

		MaxRequestsInFlight: 25,
		Timeout:             10 * time.Second,
	}))

	// Start HTTP server in the background
	executionGroup.Go(func() error {
		listenAddress := fmt.Sprintf("%s:%d", config.HTTPAddress, config.HTTPPort)
		if err := httpServer.Start(ctx, listenAddress); err != nil {
			return fmt.Errorf("HTTP server error: %w", err)
		}

		return nil
	})

	/***********************/

	mqttServer := broker.NewServer(config, isLocalDevelopment)

	// Start MQTT server in the background
	executionGroup.Go(func() error {
		slog.Info("Listening for MQTT messages", "address", config.MQTTAddress, "port", config.MQTTPort)
		if err := mqttServer.Serve(); err != nil {
			return fmt.Errorf("MQTT server error: %w", err)
		}

		return nil
	})

	/***********************/

	if err := executionGroup.Wait(); err != nil {
		log.Error("Error in execution group", "error", err)
		os.Exit(1)
	}
}
