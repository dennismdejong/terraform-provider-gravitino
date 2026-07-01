package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListRoles(metalake, resourceType, resource string) (*models.RoleListResponse, error) {
	var result models.RoleListResponse
	path := fmt.Sprintf("/metalakes/%s/%s/%s/roles",
		url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource))
	err := c.Get(path, &result)
	return &result, err
}
