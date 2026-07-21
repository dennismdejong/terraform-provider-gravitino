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

func (c *Client) ListModelVersions(metalake, catalog, schema, model string) (*models.ModelVersionListResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model))
	var result models.ModelVersionListResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) LinkModelVersion(metalake, catalog, schema, model string, req *models.ModelVersionLinkRequest) (*models.ModelVersionResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model))
	var result models.ModelVersionResponse
	if err := c.Post(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModelVersion(metalake, catalog, schema, model, version string) (*models.ModelVersionResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(version))
	var result models.ModelVersionResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateModelVersion(metalake, catalog, schema, model, version string, req *models.ModelVersionLinkRequest) (*models.ModelVersionResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(version))
	var result models.ModelVersionResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteModelVersion(metalake, catalog, schema, model, version string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(version))
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModelVersionByAlias(metalake, catalog, schema, model, alias string) (*models.ModelVersionResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/aliases/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(alias))
	var result models.ModelVersionResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteModelVersionByAlias(metalake, catalog, schema, model, alias string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/aliases/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(alias))
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateModelVersionByAlias(metalake, catalog, schema, model, alias string, req *models.ModelVersionLinkRequest) (*models.ModelVersionResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/aliases/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(alias))
	var result models.ModelVersionResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModelVersionURI(metalake, catalog, schema, model, version string) (*models.ModelVersionURIResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/%s/uri",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(version))
	var result models.ModelVersionURIResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModelVersionURIByAlias(metalake, catalog, schema, model, alias string) (*models.ModelVersionURIResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/models/%s/versions/aliases/%s/uri",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(model), url.PathEscape(alias))
	var result models.ModelVersionURIResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
