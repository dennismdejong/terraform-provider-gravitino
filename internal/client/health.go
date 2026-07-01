package client

import (
	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) GetHealth() (*models.HealthResponse, error) {
	var result models.HealthResponse
	err := c.Get("/health", &result)
	return &result, err
}

func (c *Client) GetLiveness() (*models.HealthResponse, error) {
	var result models.HealthResponse
	err := c.Get("/health/liveness", &result)
	return &result, err
}

func (c *Client) GetReadiness() (*models.HealthResponse, error) {
	var result models.HealthResponse
	err := c.Get("/health/readiness", &result)
	return &result, err
}
