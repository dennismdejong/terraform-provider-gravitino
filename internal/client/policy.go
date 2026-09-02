package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListPoliciesForObject(metalake, objType, objFullName string) (*models.NameListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/objects/%s/%s/policies",
		url.PathEscape(metalake), url.PathEscape(objType), url.PathEscape(objFullName))
	var result models.NameListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) AssociatePolicies(metalake, objType, objFullName string, req *models.PolicyAssociationRequest) (*models.NameListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/objects/%s/%s/policies",
		url.PathEscape(metalake), url.PathEscape(objType), url.PathEscape(objFullName))
	var result models.NameListResponse
	if err := c.Post(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListPolicies(metalake string) (*models.PolicyListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies", url.PathEscape(metalake))
	var result models.PolicyListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreatePolicy(metalake string, req *models.PolicyCreateRequest) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies", url.PathEscape(metalake))
	var result models.PolicyResponse
	if err := c.Post(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPolicy(metalake, name string) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.PolicyResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePolicy(metalake, name string, req *models.PolicyCreateRequest) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.PolicyResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetPolicyEnabled(metalake, name string, enabled bool) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.PolicyResponse
	body := map[string]bool{"enabled": enabled}
	if err := c.Patch(path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeletePolicy(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListObjectsForPolicy(metalake, name string) (*models.IdentifiersResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s/objects", url.PathEscape(metalake), url.PathEscape(name))
	var result models.IdentifiersResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
