package main

import (
	"log"
	"time"

	v1 "github.com/DenisPojar/go-asana-test-project/api/v1"
	"github.com/DenisPojar/go-asana-test-project/fetcher"
)

func main() {
	cfg := LoadConfig()

	client := v1.NewApiClient(cfg.APIToken)

	go fetchData(client, cfg.BaseURL, cfg.ShortInterval, "short_interval_projects.json", "short_interval_users.json")
	go fetchData(client, cfg.BaseURL, cfg.LongInterval, "long_interval_projects.json", "long_interval_users.json")

	log.Println("Service started. Fetching data... Press Ctrl+C to stop.")

	select {} // block forever
}

// fetchData periodically fetches projects and users and stores them in JSON files
func fetchData(client *v1.ApiClient, baseURL string, interval time.Duration, projectsFile string, usersFile string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		log.Printf("Fetching data for interval %s...\n", interval)

		if err := fetcher.FetchAndStore(client, projectsFile, usersFile, baseURL); err != nil {
			log.Printf("Error fetching or storing data: %v\n", err)
			continue
		}

		log.Printf("Successfully fetched and stored data to %s and %s\n", projectsFile, usersFile)
	}
}
