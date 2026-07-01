package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListFunctions(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/functions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListFunctionsDetails(metalake, catalog, schema string) ([]models.Function, error) {
	identifiers, err := c.ListFunctions(metalake, catalog, schema)
	if err != nil {
		return nil, err
	}
	functions := make([]models.Function, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetFunction(metalake, catalog, schema, id.Name)
		if err != nil {
			return nil, err
		}
		functions = append(functions, resp.Function)
	}
	return functions, nil
}

func (c *Client) GetFunction(metalake, catalog, schema, name string) (*models.FunctionResponse, error) {
	var result models.FunctionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/functions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateFunction(metalake, catalog, schema string, req *models.FunctionCreateRequest) (*models.FunctionResponse, error) {
	var result models.FunctionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/functions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateFunction(metalake, catalog, schema, name string, updates []interface{}) (*models.FunctionResponse, error) {
	var result models.FunctionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/functions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Put(path, &models.FunctionUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropFunction(metalake, catalog, schema, name string, force bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/functions/%s?force=%t",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name), force)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
