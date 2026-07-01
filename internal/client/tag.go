package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListTags(metalake string) (*models.TagNameListResponse, error) {
	var result models.TagNameListResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/tags", &result)
	return &result, err
}

func (c *Client) ListTagsDetailed(metalake string) (*models.TagListResponse, error) {
	var result models.TagListResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/tags?details=true", &result)
	return &result, err
}

func (c *Client) GetTag(metalake, name string) (*models.TagResponse, error) {
	var result models.TagResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/tags/"+url.PathEscape(name), &result)
	return &result, err
}

func (c *Client) CreateTag(metalake string, req *models.TagCreateRequest) (*models.TagResponse, error) {
	var result models.TagResponse
	err := c.Post("/metalakes/"+url.PathEscape(metalake)+"/tags", req, &result)
	return &result, err
}

func (c *Client) UpdateTag(metalake, name string, updates []interface{}) (*models.TagResponse, error) {
	var result models.TagResponse
	path := fmt.Sprintf("/metalakes/%s/tags/%s", url.PathEscape(metalake), url.PathEscape(name))
	err := c.Put(path, &models.TagUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DeleteTag(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/tags/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) ListTagsForObject(metalake, resourceType, resource string) (*models.TagListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/tags", url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource))
	var result models.TagListResponse
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) AssociateTags(metalake, resourceType, resource string, tags []string) (*models.TagListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/tags", url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource))
	var result models.TagListResponse
	err := c.Post(path, &models.TagAssociationRequest{Tags: tags}, &result)
	return &result, err
}

func (c *Client) GetTagForObject(metalake, resourceType, resource, tag string) (*models.TagResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/%s/%s/tags/%s", url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource), url.PathEscape(tag))
	var result models.TagResponse
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListObjectsForTag(metalake, tag string) (*models.BaseResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/tags/%s/objects", url.PathEscape(metalake), url.PathEscape(tag))
	var result models.BaseResponse
	err := c.Get(path, &result)
	return &result, err
}
