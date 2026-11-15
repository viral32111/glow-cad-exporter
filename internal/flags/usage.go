package flags

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

const (
	PROJECT_TITLE       = "Glow CAD Exporter"
	PROJECT_DESCRIPTION = "A Prometheus exporter for Glow CAD smart meter data, received over MQTT."
	PROJECT_URL         = "https://github.com/viral32111/glow-cad-exporter"
	PROJECT_COPYRIGHT   = "Copyright (C) 2026 viral32111, under the GNU AGPL v3.0 license."
)

const NAME_COLUMN_WIDTH = 27

type usageEntry struct {
	name        string
	description string
	fallback    string
}

func formatDuration(duration time.Duration) string {
	formatted := duration.String()

	if strings.HasSuffix(formatted, "m0s") {
		formatted = strings.TrimSuffix(formatted, "0s")
	}

	if strings.HasSuffix(formatted, "h0m") {
		formatted = strings.TrimSuffix(formatted, "0m")
	}

	return formatted
}

func PrintUsage(writer io.Writer) {
	flagEntries := []usageEntry{
		{"-v, --version", "Print version & exit.", ""},
		{"-u, --usage", "Show this usage message & exit.", ""},
	}

	environmentEntries := []usageEntry{
		{"LOG_LEVEL", "Minimum level of logs to print (debug, info, warn, error)", strings.ToLower(config.DEFAULT_LOG_LEVEL.String())},
		{"HTTP_ADDRESS", "IP address to listen for HTTP requests on", config.DEFAULT_HTTP_ADDRESS},
		{"HTTP_PORT", "TCP port to listen for HTTP requests on", strconv.Itoa(config.DEFAULT_HTTP_PORT)},
		{"METRICS_PATH", "HTTP path to expose Prometheus metrics on", config.DEFAULT_METRICS_PATH},
		{"METRICS_GO_BUILTIN", "Whether to include Go runtime & process metrics", strconv.FormatBool(config.DEFAULT_METRICS_GO_BUILTIN)},
		{"MQTT_ADDRESS", "IP address to listen for MQTT messages on", config.DEFAULT_MQTT_ADDRESS},
		{"MQTT_PORT", "TCP port to listen for MQTT messages on", strconv.Itoa(config.DEFAULT_MQTT_PORT)},
		{"MQTT_USER", "Username for authenticating MQTT connections", config.DEFAULT_MQTT_USER},
		{"MQTT_PASSWORD", "Password for authenticating MQTT connections", config.DEFAULT_MQTT_PASSWORD},
		{"MQTT_TOPIC", "Prefix for topics published by the Glow CAD", config.DEFAULT_MQTT_TOPIC},
		{"REDIS_URL", "URL to the Redis server for caching MQTT messages", config.DEFAULT_REDIS_URL},
		{"REDIS_KEY_PREFIX", "Prefix for Redis keys to avoid collisions", config.DEFAULT_REDIS_KEY_PREFIX},
		{"REDIS_STATE_TTL", "How long to cache device state MQTT messages in Redis", formatDuration(config.DEFAULT_REDIS_STATE_TTL)},
		{"REDIS_ELECTRIC_TTL", "How long to cache electricity meter MQTT messages in Redis", formatDuration(config.DEFAULT_REDIS_ELECTRIC_TTL)},
		{"REDIS_GAS_TTL", "How long to cache gas meter MQTT messages in Redis", formatDuration(config.DEFAULT_REDIS_GAS_TTL)},
	}

	for index := range environmentEntries {
		environmentEntries[index].name = config.ENVIRONMENT_VARIABLE_PREFIX + environmentEntries[index].name
	}

	columnWidth := NAME_COLUMN_WIDTH
	for _, entry := range append(flagEntries, environmentEntries...) {
		columnWidth = max(columnWidth, len(entry.name)+1)
	}

	fmt.Fprintf(writer, "%s - %s\n", PROJECT_TITLE, PROJECT_DESCRIPTION)
	fmt.Fprintf(writer, "  %s\n", PROJECT_URL)
	fmt.Fprintf(writer, "\n")

	fmt.Fprintf(writer, "Usage:\n")
	fmt.Fprintf(writer, "  %s [flags]\n", os.Args[0])
	fmt.Fprintf(writer, "\n")

	fmt.Fprintf(writer, "Flags:\n")
	for _, entry := range flagEntries {
		fmt.Fprintf(writer, "  %-*s%s\n", columnWidth, entry.name, entry.description)
	}
	fmt.Fprintf(writer, "\n")

	fmt.Fprintf(writer, "Environment variables:\n")
	for _, entry := range environmentEntries {
		fmt.Fprintf(writer, "  %-*s%s (default: '%s')\n", columnWidth, entry.name, entry.description, entry.fallback)
	}
	fmt.Fprintf(writer, "\n")

	fmt.Fprintf(writer, "%s\n", PROJECT_COPYRIGHT)
}
