package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) GetCredentials(metalake, resourceType, resource string) (*models.CredentialResponse, error) {
	var result models.CredentialResponse
	path := fmt.Sprintf("/metalakes/%s/objects/%s/%s/credentials",
		url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource))
	err := c.Get(path, &result)
	return &result, err
}
