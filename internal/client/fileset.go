package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListFilesets(metalake, catalog, schema string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/filesets",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) GetFileset(metalake, catalog, schema, name string) (*models.FilesetResponse, error) {
	var result models.FilesetResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/filesets/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreateFileset(metalake, catalog, schema string, req *models.FilesetCreateRequest) (*models.FilesetResponse, error) {
	var result models.FilesetResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/filesets",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdateFileset(metalake, catalog, schema, name string, updates []interface{}) (*models.FilesetResponse, error) {
	var result models.FilesetResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/filesets/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	err := c.Put(path, &models.FilesetUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropFileset(metalake, catalog, schema, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/filesets/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) ListFilesetsDetails(metalake, catalog, schema string) ([]models.Fileset, error) {
	identifiers, err := c.ListFilesets(metalake, catalog, schema)
	if err != nil {
		return nil, err
	}
	filesets := make([]models.Fileset, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetFileset(metalake, catalog, schema, id.Name)
		if err != nil {
			return nil, err
		}
		filesets = append(filesets, resp.Fileset)
	}
	return filesets, nil
}
