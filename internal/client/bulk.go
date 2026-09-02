package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) BulkAddUsers(metalake string, names []string) (*models.BulkUserResponse, error) {
	users := make([]models.UserCreateRequest, 0, len(names))
	for _, n := range names {
		users = append(users, models.UserCreateRequest{Name: n})
	}
	path := fmt.Sprintf("/bulk/metalakes/%s/users/add", url.PathEscape(metalake))
	var result models.BulkUserResponse
	if err := c.Post(path, &models.BulkUserAddRequest{Users: users}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BulkRemoveUsers(metalake string, names []string) (*models.BulkRemoveResponse, error) {
	path := fmt.Sprintf("/bulk/metalakes/%s/users/remove", url.PathEscape(metalake))
	var result models.BulkRemoveResponse
	if err := c.Post(path, &models.BulkRemoveRequest{Names: names}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BulkAddGroups(metalake string, names []string) (*models.BulkGroupResponse, error) {
	groups := make([]models.GroupCreateRequest, 0, len(names))
	for _, n := range names {
		groups = append(groups, models.GroupCreateRequest{Name: n})
	}
	path := fmt.Sprintf("/bulk/metalakes/%s/groups/add", url.PathEscape(metalake))
	var result models.BulkGroupResponse
	if err := c.Post(path, &models.BulkGroupAddRequest{Groups: groups}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BulkRemoveGroups(metalake string, names []string) (*models.BulkRemoveResponse, error) {
	path := fmt.Sprintf("/bulk/metalakes/%s/groups/remove", url.PathEscape(metalake))
	var result models.BulkRemoveResponse
	if err := c.Post(path, &models.BulkRemoveRequest{Names: names}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
