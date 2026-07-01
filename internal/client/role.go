package client

import (
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListRoles(metalake, resourceType, resource string) (*models.RoleListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/objects/%s/%s/roles", metalake, resourceType, resource)
	var result models.RoleListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateRole(metalake string, req *models.RoleCreateRequest) (*models.RoleResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/roles", metalake)
	var result models.RoleResponse
	if err := c.Post(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetRole(metalake, name string) (*models.RoleResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/roles/%s", metalake, name)
	var result models.RoleResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteRole(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/roles/%s", metalake, name)
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListAllRoles(metalake string) (*models.NameListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/roles", metalake)
	var result models.NameListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GrantPrivilegeToRole(metalake, role, objectType, objectFullName string, privileges []models.Privilege) (*models.RoleResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/roles/%s/%s/%s/grant", metalake, role, objectType, objectFullName)
	var result models.RoleResponse
	if err := c.Put(path, &models.PrivilegesRequest{Privileges: privileges}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RevokePrivilegeFromRole(metalake, role, objectType, objectFullName string, privileges []models.Privilege) (*models.RoleResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/roles/%s/%s/%s/revoke", metalake, role, objectType, objectFullName)
	var result models.RoleResponse
	if err := c.Put(path, &models.PrivilegesRequest{Privileges: privileges}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) OverrideRolePrivileges(metalake, role, objectType, objectFullName string, privileges []models.Privilege) (*models.RoleResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/permissions/roles/%s/%s/%s/override", metalake, role, objectType, objectFullName)
	var result models.RoleResponse
	if err := c.Put(path, &models.PrivilegesRequest{Privileges: privileges}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
