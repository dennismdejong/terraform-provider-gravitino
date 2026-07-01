package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListModels(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListModelsDetails(metalake, catalog, schema string) ([]models.Model, error) {
	identifiers, err := c.ListModels(metalake, catalog, schema)
	if err != nil {
		return nil, err
	}
	mods := make([]models.Model, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetModel(metalake, catalog, schema, id.Name)
		if err != nil {
			return nil, err
		}
		mods = append(mods, resp.Model)
	}
	return mods, nil
}

func (c *Client) GetModel(metalake, catalog, schema, model string) (*models.ModelResponse, error) {
	var result models.ModelResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateModel(metalake, catalog, schema string, req *models.ModelCreateRequest) (*models.ModelResponse, error) {
	var result models.ModelResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateModel(metalake, catalog, schema, model string, updates []interface{}) (*models.ModelResponse, error) {
	var result models.ModelResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model))
	err := c.Put(path, &models.ModelUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropModel(metalake, catalog, schema, model string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
