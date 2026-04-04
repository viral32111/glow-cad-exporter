package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/viral32111/glow-cad-exporter/internal/config"
	"github.com/viral32111/glow-cad-exporter/internal/mqtt"
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

	signals := make(chan os.Signal, 1)
	isDone := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals

		slog.Info("Caught termination signal, stopping...")
		isDone <- true
	}()

	/***********************/

	mqttServer := mochi.New(&mochi.Options{})

	// Only output warnings & errors from the MQTT server itself
	mqttLogOptions := &slog.HandlerOptions{Level: slog.LevelDebug}
	if isLocalDevelopment {
		mqttServer.Log = slog.New(slog.NewTextHandler(os.Stdout, mqttLogOptions))
	} else {
		mqttServer.Log = slog.New(slog.NewJSONHandler(os.Stdout, mqttLogOptions))
	}

	// Require authentication if username/password is set, otherwise allow anonymous connections
	if (config.MQTTUser != "") && (config.MQTTPassword != "") {
		// https://github.com/mochi-mqtt/server?tab=readme-ov-file#access-control
		if err := mqttServer.AddHook(new(auth.Hook), &auth.Options{
			Ledger: &auth.Ledger{
				Auth: auth.AuthRules{
					{
						Username: auth.RString(config.MQTTUser),
						Password: auth.RString(config.MQTTPassword),
						Allow:    true,
					},
				},
			},
		}); err != nil {
			slog.Error("Unable to add MQTT auth hook", "error", err)
			os.Exit(1)
		}
	} else {
		if err := mqttServer.AddHook(new(auth.AllowHook), nil); err != nil {
			slog.Error("Unable to add MQTT anonymous hook", "error", err)
			os.Exit(1)
		}
	}

	if err := mqttServer.AddHook(&mqtt.Hooks{
		Config: config,
	}, nil); err != nil {
		slog.Error("Unable to add MQTT glow hook", "error", err)
		os.Exit(1)
	}

	tcpListener := listeners.NewTCP(listeners.Config{
		ID:      "tcp0",
		Address: fmt.Sprintf("%s:%d", config.MQTTAddress, config.MQTTPort),
	})
	if err := mqttServer.AddListener(tcpListener); err != nil {
		slog.Error("Unable to add TCP listener to MQTT server", "error", err)
		os.Exit(1)
	}

	// Start MQTT server in the background
	go func() {
		slog.Info("Listening for MQTT messages", "address", config.MQTTAddress, "port", config.MQTTPort)
		if err := mqttServer.Serve(); err != nil {
			slog.Error("MQTT server error", "error", err)
		}
	}()

	// Wait for termination signal
	<-isDone
}
