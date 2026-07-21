package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListPolicies(metalake, resourceType, resource string) (*models.PolicyListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/policies",
		url.PathEscape(metalake),
		url.PathEscape(resourceType),
		url.PathEscape(resource),
	)
	var result models.PolicyListResponse
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) GetPolicy(metalake, resourceType, resource, name string) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/policies/%s",
		url.PathEscape(metalake),
		url.PathEscape(resourceType),
		url.PathEscape(resource),
		url.PathEscape(name),
	)
	var result models.PolicyResponse
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreatePolicy(metalake, resourceType, resource string, req *models.PolicyCreateRequest) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/policies",
		url.PathEscape(metalake),
		url.PathEscape(resourceType),
		url.PathEscape(resource),
	)
	var result models.PolicyResponse
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdatePolicy(metalake, resourceType, resource, name string, updates []interface{}) (*models.PolicyResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/policies/%s",
		url.PathEscape(metalake),
		url.PathEscape(resourceType),
		url.PathEscape(resource),
		url.PathEscape(name),
	)
	var result models.PolicyResponse
	err := c.Put(path, &models.PolicyUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DeletePolicy(metalake, resourceType, resource, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/policies/%s",
		url.PathEscape(metalake),
		url.PathEscape(resourceType),
		url.PathEscape(resource),
		url.PathEscape(name),
	)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) AssociatePoliciesForObject(metalake, objType, objFullName string, policyNames []string) (*models.BaseResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/objects/%s/%s/policies",
		url.PathEscape(metalake), url.PathEscape(objType), url.PathEscape(objFullName))
	var result models.BaseResponse
	body := map[string][]string{"policyNames": policyNames}
	if err := c.Post(path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetPolicy(metalake, policyName string, enabled bool) (*models.BaseResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s",
		url.PathEscape(metalake), url.PathEscape(policyName))
	var result models.BaseResponse
	body := map[string]bool{"enabled": enabled}
	if err := c.Patch(path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListMetadataObjectsForPolicy(metalake, policyName string) (*models.IdentifiersResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/policies/%s/objects",
		url.PathEscape(metalake), url.PathEscape(policyName))
	var result models.IdentifiersResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
