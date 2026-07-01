package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListTopics(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/topics",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListTopicsDetails(metalake, catalog, schema string) ([]models.Topic, error) {
	identifiers, err := c.ListTopics(metalake, catalog, schema)
	if err != nil {
		return nil, err
	}
	topics := make([]models.Topic, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetTopic(metalake, catalog, schema, id.Name)
		if err != nil {
			return nil, err
		}
		topics = append(topics, resp.Topic)
	}
	return topics, nil
}

func (c *Client) GetTopic(metalake, catalog, schema, name string) (*models.TopicResponse, error) {
	var result models.TopicResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/topics/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateTopic(metalake, catalog, schema string, req *models.TopicCreateRequest) (*models.TopicResponse, error) {
	var result models.TopicResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/topics",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateTopic(metalake, catalog, schema, name string, updates []interface{}) (*models.TopicResponse, error) {
	var result models.TopicResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/topics/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Put(path, &models.TopicUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropTopic(metalake, catalog, schema, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/topics/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
