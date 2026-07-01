package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListViews(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/views",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListViewsDetails(metalake, catalog, schema string) ([]models.View, error) {
	identifiers, err := c.ListViews(metalake, catalog, schema)
	if err != nil {
		return nil, err
	}
	views := make([]models.View, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetView(metalake, catalog, schema, id.Name)
		if err != nil {
			return nil, err
		}
		views = append(views, resp.View)
	}
	return views, nil
}

func (c *Client) GetView(metalake, catalog, schema, name string) (*models.ViewResponse, error) {
	var result models.ViewResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/views/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateView(metalake, catalog, schema string, req *models.ViewCreateRequest) (*models.ViewResponse, error) {
	var result models.ViewResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/views",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateView(metalake, catalog, schema, name string, updates []interface{}) (*models.ViewResponse, error) {
	var result models.ViewResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/views/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Put(path, &models.ViewUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropView(metalake, catalog, schema, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/views/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}
