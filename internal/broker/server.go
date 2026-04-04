package broker

import (
	"fmt"
	"log/slog"
	"os"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

func NewServer(config config.Config, useTextLogger bool) *mqtt.Server {
	server := mqtt.New(nil)

	// Only output warnings & errors from the MQTT server itself
	logOptions := &slog.HandlerOptions{Level: slog.LevelWarn}
	if useTextLogger {
		server.Log = slog.New(slog.NewTextHandler(os.Stdout, logOptions))
	} else {
		server.Log = slog.New(slog.NewJSONHandler(os.Stdout, logOptions))
	}

	// Require authentication if username/password is set, otherwise allow anonymous connections
	if (config.MQTTUser != "") && (config.MQTTPassword != "") {
		// https://github.com/mochi-mqtt/server?tab=readme-ov-file#access-control
		if err := server.AddHook(new(auth.Hook), &auth.Options{
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
		if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
			slog.Error("Unable to add MQTT anonymous hook", "error", err)
			os.Exit(1)
		}
	}

	if err := server.AddHook(&Hooks{
		Config: config,
	}, nil); err != nil {
		slog.Error("Unable to add MQTT glow hook", "error", err)
		os.Exit(1)
	}

	tcpListener := listeners.NewTCP(listeners.Config{
		ID:      "tcp0",
		Address: fmt.Sprintf("%s:%d", config.MQTTAddress, config.MQTTPort),
	})
	if err := server.AddListener(tcpListener); err != nil {
		slog.Error("Unable to add TCP listener to MQTT server", "error", err)
		os.Exit(1)
	}

	return server
}
