package v1

import (
	"encoding/json"
	"fmt"
)

// FetchProjects fetches projects using fetchWithRetry
func FetchProjects(c *ApiClient, baseURL string) ([]Project, error) {
	url := fmt.Sprintf("%s/projects", baseURL)

	body, err := c.fetchWithRetry(url, 5)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []Project `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return resp.Data, nil
}
