package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListSchemas(metalake, catalog string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas", url.PathEscape(metalake), url.PathEscape(catalog))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) GetSchema(metalake, catalog, schema string) (*models.SchemaResponse, error) {
	var result models.SchemaResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateSchema(metalake, catalog string, req *models.SchemaCreateRequest) (*models.SchemaResponse, error) {
	var result models.SchemaResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas", url.PathEscape(metalake), url.PathEscape(catalog))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateSchema(metalake, catalog, schema string, updates []interface{}) (*models.SchemaResponse, error) {
	var result models.SchemaResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Put(path, &models.SchemaUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropSchema(metalake, catalog, schema string, force bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s?force=%t",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), force)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) ListSchemasDetails(metalake, catalog string) ([]models.Schema, error) {
	identifiers, err := c.ListSchemas(metalake, catalog)
	if err != nil {
		return nil, err
	}
	schemas := make([]models.Schema, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetSchema(metalake, catalog, id.Name)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, resp.Schema)
	}
	return schemas, nil
}
