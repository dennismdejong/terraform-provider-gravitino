package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListTables(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) GetTable(metalake, catalog, schema, table string) (*models.TableResponse, error) {
	var result models.TableResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateTable(metalake, catalog, schema string, req *models.TableCreateRequest) (*models.TableResponse, error) {
	var result models.TableResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateTable(metalake, catalog, schema, table string, updates []interface{}) (*models.TableResponse, error) {
	var result models.TableResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table))
	err := c.Put(path, &models.TableUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropTable(metalake, catalog, schema, table string, force, purge bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s?force=%t&purge=%t",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table), force, purge)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
