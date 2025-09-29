package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIToken      string
	BaseURL       string
	ShortInterval time.Duration
	LongInterval  time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		log.Fatal("API_TOKEN is not set in the environment")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL is not set in the environment")
	}

	shortIntervalStr := os.Getenv("SHORT_INTERVAL")
	if shortIntervalStr == "" {
		log.Fatal("SHORT_INTERVAL is not set in the environment")
	}
	shortInterval, err := time.ParseDuration(shortIntervalStr)
	if err != nil {
		log.Fatalf("invalid SHORT_INTERVAL: %v", err)
	}

	longIntervalStr := os.Getenv("LONG_INTERVAL")
	if longIntervalStr == "" {
		log.Fatal("LONG_INTERVAL is not set in the environment")
	}
	longInterval, err := time.ParseDuration(longIntervalStr)
	if err != nil {
		log.Fatalf("invalid LONG_INTERVAL: %v", err)
	}

	return &Config{
		APIToken:      apiToken,
		BaseURL:       baseURL,
		ShortInterval: shortInterval,
		LongInterval:  longInterval,
	}

}
