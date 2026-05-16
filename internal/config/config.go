package config

import (
	"os"
	"strconv"
)

// Runtime configuration from environment variables
type Config struct {
	Port        string
	MaxFileSize int64
	RateLimit   float64
	RateBurst   int
}

func Load() *Config {
	return &Config{
		Port:        getStr("PORT", "3000"),
		MaxFileSize: getInt64("MAX_FILE_SIZE", 10<<20),
		RateLimit:   getFloat64("RATE_LIMIT", 5),
		RateBurst:   getInt("RATE_BURST", 10),
	}
}

func getStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getFloat64(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
