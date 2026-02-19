package config

import (
	"os"
	"time"
)

type Config struct {
	ServerPort      string
	AnalyzerURL     string
	AnalysisTimeout time.Duration
}

type ConfigDB struct {
	URL string
}

func NewConfig() *Config {
	timeoutStr := getEnv("ANALYSIS_TIMEOUT", "5s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = 5 * time.Second
	}

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		AnalyzerURL:     getEnv("ANALYZER_URL", "8081"),
		AnalysisTimeout: timeout,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
