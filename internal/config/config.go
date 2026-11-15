package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"
	"unicode"
)

const (
	DEFAULT_LOG_LEVEL = slog.LevelInfo

	DEFAULT_HTTP_ADDRESS = "127.0.0.1"
	DEFAULT_HTTP_PORT    = 5000

	DEFAULT_METRICS_PATH       = "/metrics"
	DEFAULT_METRICS_GO_BUILTIN = false

	DEFAULT_MQTT_ADDRESS  = "127.0.0.1"
	DEFAULT_MQTT_PORT     = 1883
	DEFAULT_MQTT_USER     = ""
	DEFAULT_MQTT_PASSWORD = ""
	DEFAULT_MQTT_TOPIC    = "glow"

	DEFAULT_REDIS_URL          = "redis://default@127.0.0.1:6379/0"
	DEFAULT_REDIS_KEY_PREFIX   = "glow"
	DEFAULT_REDIS_STATE_TTL    = 24 * time.Hour
	DEFAULT_REDIS_ELECTRIC_TTL = 730 * time.Hour
	DEFAULT_REDIS_GAS_TTL      = 730 * time.Hour
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
	logLevel, err := parseEnvLogLevel("LOG_LEVEL", DEFAULT_LOG_LEVEL)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "LOG_LEVEL", err)
	}

	httpPort, err := parseEnvNumber("HTTP_PORT", DEFAULT_HTTP_PORT)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "HTTP_PORT", err)
	}

	if httpPort < 1 || httpPort > 65535 {
		return Config{}, fmt.Errorf("Invalid value for environment variable '%s%s': Port %d must be between 1 and 65535", ENVIRONMENT_VARIABLE_PREFIX, "HTTP_PORT", httpPort)
	}

	metricsGoBuiltin, err := parseEnvBool("METRICS_GO_BUILTIN", DEFAULT_METRICS_GO_BUILTIN)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "METRICS_GO_BUILTIN", err)
	}

	metricsPath := parseEnvString("METRICS_PATH", DEFAULT_METRICS_PATH)

	if !strings.HasPrefix(metricsPath, "/") {
		return Config{}, fmt.Errorf("Invalid value for environment variable '%s%s': Path %q must start with a slash", ENVIRONMENT_VARIABLE_PREFIX, "METRICS_PATH", metricsPath)
	}

	if len(metricsPath) > 1 && strings.HasSuffix(metricsPath, "/") {
		return Config{}, fmt.Errorf("Invalid value for environment variable '%s%s': Path %q must not end with a slash", ENVIRONMENT_VARIABLE_PREFIX, "METRICS_PATH", metricsPath)
	}

	if strings.ContainsFunc(metricsPath, unicode.IsSpace) {
		return Config{}, fmt.Errorf("Invalid value for environment variable '%s%s': Path %q must not contain whitespace", ENVIRONMENT_VARIABLE_PREFIX, "METRICS_PATH", metricsPath)
	}

	mqttPort, err := parseEnvNumber("MQTT_PORT", DEFAULT_MQTT_PORT)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "MQTT_PORT", err)
	}

	if mqttPort < 1 || mqttPort > 65535 {
		return Config{}, fmt.Errorf("Invalid value for environment variable '%s%s': Port %d must be between 1 and 65535", ENVIRONMENT_VARIABLE_PREFIX, "MQTT_PORT", mqttPort)
	}

	redisStateTTL, err := parseEnvDuration("REDIS_STATE_TTL", DEFAULT_REDIS_STATE_TTL)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_STATE_TTL", err)
	}

	redisElectricTTL, err := parseEnvDuration("REDIS_ELECTRIC_TTL", DEFAULT_REDIS_ELECTRIC_TTL)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_ELECTRIC_TTL", err)
	}

	redisGasTTL, err := parseEnvDuration("REDIS_GAS_TTL", DEFAULT_REDIS_GAS_TTL)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to parse environment variable '%s%s': %w", ENVIRONMENT_VARIABLE_PREFIX, "REDIS_GAS_TTL", err)
	}

	config := Config{
		LogLevel: logLevel,

		HTTPAddress: parseEnvString("HTTP_ADDRESS", DEFAULT_HTTP_ADDRESS),
		HTTPPort:    httpPort,

		MetricsPath:      metricsPath,
		MetricsGoBuiltin: metricsGoBuiltin,

		MQTTAddress:  parseEnvString("MQTT_ADDRESS", DEFAULT_MQTT_ADDRESS),
		MQTTPort:     mqttPort,
		MQTTUser:     parseEnvString("MQTT_USER", DEFAULT_MQTT_USER),
		MQTTPassword: parseEnvString("MQTT_PASSWORD", DEFAULT_MQTT_PASSWORD),
		MQTTTopic:    parseEnvString("MQTT_TOPIC", DEFAULT_MQTT_TOPIC),

		RedisURL:         parseEnvString("REDIS_URL", DEFAULT_REDIS_URL),
		RedisKeyPrefix:   parseEnvString("REDIS_KEY_PREFIX", DEFAULT_REDIS_KEY_PREFIX),
		RedisStateTTL:    redisStateTTL,
		RedisElectricTTL: redisElectricTTL,
		RedisGasTTL:      redisGasTTL,
	}

	// The Redis URL may embed a password, so it is redacted before being logged
	redisURL := "<unparseable>"
	if parsedRedisURL, err := url.Parse(config.RedisURL); err == nil {
		redisURL = parsedRedisURL.Redacted()
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

		"redis_url", redisURL,
		"redis_key_prefix", config.RedisKeyPrefix,
		"redis_state_ttl", config.RedisStateTTL,
		"redis_electric_ttl", config.RedisElectricTTL,
		"redis_gas_ttl", config.RedisGasTTL,
	)

	return config, nil
}
