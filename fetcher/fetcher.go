package fetcher

import (
	"fmt"

	v1 "github.com/DenisPojar/go-asana-test-project/api/v1"
	"github.com/DenisPojar/go-asana-test-project/storage"
)

// FetchAndStore fetches projects and users, then stores them in JSON files
func FetchAndStore(client *v1.ApiClient, projectsFile string, usersFile string, baseURL string) error {
	// 1. Fetch Projects
	projects, err := v1.FetchProjects(client, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	// 2. Save Projects to file
	if err := storage.SaveJSON(projectsFile, projects); err != nil {
		return fmt.Errorf("failed to save projects: %w", err)
	}

	// 3. Fetch Users
	users, err := v1.FetchUsers(client, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	// 4. Save Users to file
	if err := storage.SaveJSON(usersFile, users); err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}

	return nil
}
