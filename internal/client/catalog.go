package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListCatalogs(metalake string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/catalogs", &result)
	return &result, err
}

func (c *Client) ListCatalogsDetails(metalake string) (*models.CatalogInfoListResponse, error) {
	var result models.CatalogInfoListResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/catalogs?details=true", &result)
	return &result, err
}

func (c *Client) GetCatalog(metalake, name string) (*models.CatalogResponse, error) {
	var result models.CatalogResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/catalogs/"+url.PathEscape(name), &result)
	return &result, err
}

func (c *Client) CreateCatalog(metalake string, req *models.CatalogCreateRequest) (*models.CatalogResponse, error) {
	var result models.CatalogResponse
	err := c.Post("/metalakes/"+url.PathEscape(metalake)+"/catalogs", req, &result)
	return &result, err
}

func (c *Client) UpdateCatalog(metalake, name string, updates []interface{}) (*models.CatalogResponse, error) {
	var result models.CatalogResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s", url.PathEscape(metalake), url.PathEscape(name))
	err := c.Put(path, &models.CatalogUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) SetCatalogInUse(metalake, name string) (*models.BaseResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/set", url.PathEscape(metalake), url.PathEscape(name))
	var result models.BaseResponse
	err := c.Put(path, nil, &result)
	return &result, err
}

func (c *Client) DropCatalog(metalake, name string, force bool) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s?force=%t", url.PathEscape(metalake), url.PathEscape(name), force)
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) TestCatalogConnection(metalake, catalogName string, testReq interface{}) (*models.BaseResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/test",
		url.PathEscape(metalake), url.PathEscape(catalogName))
	var result models.BaseResponse
	if err := c.Post(path, testReq, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
