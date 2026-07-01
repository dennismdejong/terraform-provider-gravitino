package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListMetalakes() (*models.MetalakeListResponse, error) {
	var result models.MetalakeListResponse
	err := c.Get("/metalakes", &result)
	return &result, err
}

func (c *Client) GetMetalake(name string) (*models.MetalakeResponse, error) {
	var result models.MetalakeResponse
	err := c.Get("/metalakes/"+url.PathEscape(name), &result)
	return &result, err
}

func (c *Client) CreateMetalake(req *models.MetalakeCreateRequest) (*models.MetalakeResponse, error) {
	var result models.MetalakeResponse
	err := c.Post("/metalakes", req, &result)
	return &result, err
}

func (c *Client) UpdateMetalake(name string, updates []interface{}) (*models.MetalakeResponse, error) {
	var result models.MetalakeResponse
	err := c.Put("/metalakes/"+url.PathEscape(name), &models.MetalakeUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropMetalake(name string, force bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s?force=%t", url.PathEscape(name), force)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
