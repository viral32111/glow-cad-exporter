package config

import (
	"fmt"
	"log/slog"
	"time"
)

type Config struct {
	LogLevel slog.Level

	HTTPAddress string
	HTTPPort    int

	MetricsPath      string
	MetricsGoBuiltin bool

	MQTTAddress  string
	MQTTPort     int
	MQTTUser     string
	MQTTPassword string
	MQTTTopic    string

	RedisURL         string
	RedisKeyPrefix   string
	RedisStateTTL    time.Duration
	RedisElectricTTL time.Duration
	RedisGasTTL      time.Duration
}

func GetConfig() (Config, error) {
	logLevel, err := parseEnvLogLevel("LOG_LEVEL", slog.LevelInfo)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "LOG_LEVEL", err)
	}

	httpPort, err := parseEnvNumber("HTTP_PORT", 8080)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "HTTP_PORT", err)
	}

	metricsGoBuiltin, err := parseEnvBool("METRICS_GO_BUILTIN", false)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "METRICS_GO_BUILTIN", err)
	}

	mqttPort, err := parseEnvNumber("MQTT_PORT", 1883)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "MQTT_PORT", err)
	}

	redisStateTTL, err := parseEnvDuration("REDIS_STATE_TTL", 24*time.Hour)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_STATE_TTL", err)
	}

	redisElectricTTL, err := parseEnvDuration("REDIS_ELECTRIC_TTL", 730*time.Hour)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_ELECTRIC_TTL", err)
	}

	redisGasTTL, err := parseEnvDuration("REDIS_GAS_TTL", 730*time.Hour)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_GAS_TTL", err)
	}

	config := Config{
		LogLevel: logLevel,

		HTTPAddress: parseEnvString("HTTP_ADDRESS", "127.0.0.1"),
		HTTPPort:    httpPort,

		MetricsPath:      parseEnvString("METRICS_PATH", "/metrics"),
		MetricsGoBuiltin: metricsGoBuiltin,

		MQTTAddress:  parseEnvString("MQTT_ADDRESS", "127.0.0.1"),
		MQTTPort:     mqttPort,
		MQTTUser:     parseEnvString("MQTT_USER", ""),
		MQTTPassword: parseEnvString("MQTT_PASSWORD", ""),
		MQTTTopic:    parseEnvString("MQTT_TOPIC", "glow"),

		RedisURL:         parseEnvString("REDIS_URL", "redis://default@127.0.0.1:6379/0"),
		RedisKeyPrefix:   parseEnvString("REDIS_KEY_PREFIX", "glow"),
		RedisStateTTL:    redisStateTTL,
		RedisElectricTTL: redisElectricTTL,
		RedisGasTTL:      redisGasTTL,
	}

	slog.Debug("Loaded configuration",
		"log_level", config.LogLevel,

		"http_address", config.HTTPAddress,
		"http_port", config.HTTPPort,

		"metrics_path", config.MetricsPath,
		"metrics_go_builtin", config.MetricsGoBuiltin,

		"mqtt_address", config.MQTTAddress,
		"mqtt_port", config.MQTTPort,
		// "mqtt_user", config.MQTTUser,
		// "mqtt_password", config.MQTTPassword,
		"mqtt_topic", config.MQTTTopic,

		"redis_url", config.RedisURL,
		"redis_key_prefix", config.RedisKeyPrefix,
		"redis_state_ttl", config.RedisStateTTL,
		"redis_electric_ttl", config.RedisElectricTTL,
		"redis_gas_ttl", config.RedisGasTTL,
	)

	return config, nil
}
