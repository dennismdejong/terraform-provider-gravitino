package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) AddIdpUser(req *models.IdpAddUserRequest) (*models.IdpUserResponse, error) {
	var result models.IdpUserResponse
	if err := c.Post("/idp/users", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetIdpUser(user string) (*models.IdpUserResponse, error) {
	path := fmt.Sprintf("/idp/users/%s", url.PathEscape(user))
	var result models.IdpUserResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateIdpUser(user string, req *models.IdpUpdateUserRequest) (*models.IdpUserResponse, error) {
	path := fmt.Sprintf("/idp/users/%s", url.PathEscape(user))
	var result models.IdpUserResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RemoveIdpUser(user string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/idp/users/%s", url.PathEscape(user))
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) AddIdpGroup(req *models.IdpAddGroupRequest) (*models.IdpGroupResponse, error) {
	var result models.IdpGroupResponse
	if err := c.Post("/idp/groups", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetIdpGroup(group string) (*models.IdpGroupResponse, error) {
	path := fmt.Sprintf("/idp/groups/%s", url.PathEscape(group))
	var result models.IdpGroupResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RemoveIdpGroup(group string, force bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/idp/groups/%s?force=%t", url.PathEscape(group), force)
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ChangeIdpGroupMembership(group string, req *models.IdpGroupMembershipChangeRequest) (*models.IdpGroupResponse, error) {
	path := fmt.Sprintf("/idp/groups/%s/users", url.PathEscape(group))
	var result models.IdpGroupResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
