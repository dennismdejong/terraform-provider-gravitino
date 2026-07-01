package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) GetAuthenticatedPrincipal() (*models.PrincipalResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/principal", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", contentType)
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("request failed with status %d: %s (%s)", resp.StatusCode, errResp.Message, errResp.Type)
	}

	var result models.PrincipalResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
