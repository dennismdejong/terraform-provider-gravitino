package client

import (
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) AddGroup(metalake string, name string) (*models.GroupResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/groups", metalake)
	var result models.GroupResponse
	if err := c.Post(path, &models.GroupCreateRequest{Name: name}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetGroup(metalake, name string) (*models.GroupResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/groups/%s", metalake, name)
	var result models.GroupResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListGroups(metalake string) (*models.NameListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/groups", metalake)
	var result models.NameListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RemoveGroup(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/groups/%s", metalake, name)
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GrantRolesToGroup(metalake, group string, roleNames []string) (*models.GroupResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/groups/%s/grant", metalake, group)
	var result models.GroupResponse
	if err := c.Put(path, &models.GrantRevokeRequest{RoleNames: roleNames}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RevokeRolesFromGroup(metalake, group string, roleNames []string) (*models.GroupResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/groups/%s/revoke", metalake, group)
	var result models.GroupResponse
	if err := c.Put(path, &models.GrantRevokeRequest{RoleNames: roleNames}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
