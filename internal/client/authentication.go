package client

import (
	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) GetAuthenticatedPrincipal() (*models.PrincipalResponse, error) {
	var result models.PrincipalResponse
	err := c.Get("/authn/me", &result)
	return &result, err
}
