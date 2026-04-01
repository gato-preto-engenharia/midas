package midas

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is midas configuration abstraction
type Config struct {
	vars map[string]string
}

// NewConfigFromEnv Create a new [Config] from the os environment variables ([os.Environ])
func NewConfigFromEnv() Config {
	cfg := Config{
		vars: make(map[string]string),
	}

	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		cfg.vars[key] = value
	}

	return cfg
}

// NewConfigFromJson Create a new [Config] from a JSON byte array
func NewConfigFromJson(jsonBytes []byte) Config {
	cfg := Config{
		vars: make(map[string]string),
	}

	if err := json.Unmarshal(jsonBytes, &cfg.vars); err != nil {
		return Config{}
	}

	return cfg
}

// Get Fetch a config variable and, if not present, use the fallback
func (c Config) Get(key, fallback string) string {
	value, exists := c.vars[key]
	if !exists {
		return fallback
	}

	return value
}

// GetInt Has the same behaviour as [Config.Get] but converts the result to int
// using [strconv.Atoi]
func (c Config) GetInt(key string, fallback int) int {
	value, exists := c.vars[key]
	if !exists {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

// GetBool Has the same behaviour as [Config.Get] but converts the result to bool
// returning true if the var value is "true", and [false] if "false"
func (c Config) GetBool(key string, fallback bool) bool {
	value, exists := c.vars[key]
	if !exists {
		return fallback
	}

	return strings.ToLower(value) == "true"
}

// GetDuration Has the same behaviour as [Config.Get] but converts the result to [time.Duration]
// using [time.ParseDuration]
func (c Config) GetDuration(key string, fallback time.Duration) time.Duration {
	value, exists := c.vars[key]
	if !exists {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

// MergeConfig Merge multiple configs into a new one, doing the overwrites in order
func MergeConfig(configs ...Config) Config {
	cfg := Config{}

	for _, c := range configs {
		for k, v := range c.vars {
			cfg.vars[k] = v
		}
	}

	return cfg
}
