package client

import (
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) GetOwner(metalake, objectType, objectFullName string) (*models.OwnerResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/owners/%s/%s", metalake, objectType, objectFullName)
	var result models.OwnerResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetOwner(metalake, objectType, objectFullName string, req *models.SetOwnerRequest) (*models.SetOwnerResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/owners/%s/%s", metalake, objectType, objectFullName)
	var result models.SetOwnerResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
