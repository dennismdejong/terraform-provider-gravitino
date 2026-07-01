package client

import (
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) AddUser(metalake string, name string) (*models.UserResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/users", metalake)
	var result models.UserResponse
	if err := c.Post(path, &models.UserCreateRequest{Name: name}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetUser(metalake, name string) (*models.UserResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/users/%s", metalake, name)
	var result models.UserResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListUsers(metalake string) (*models.NameListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/users", metalake)
	var result models.NameListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RemoveUser(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/users/%s", metalake, name)
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GrantRolesToUser(metalake, user string, roleNames []string) (*models.UserResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/users/%s/grant", metalake, user)
	var result models.UserResponse
	if err := c.Put(path, &models.GrantRevokeRequest{RoleNames: roleNames}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RevokeRolesFromUser(metalake, user string, roleNames []string) (*models.UserResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/users/%s/revoke", metalake, user)
	var result models.UserResponse
	if err := c.Put(path, &models.GrantRevokeRequest{RoleNames: roleNames}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
