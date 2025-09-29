package v1

import (
	"encoding/json"
	"fmt"
)

func FetchUsers(c *ApiClient, baseURL string) ([]User, error) {
	url := fmt.Sprintf("%s/users", baseURL)

	body, err := c.fetchWithRetry(url, 5)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []User `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return resp.Data, nil
}
