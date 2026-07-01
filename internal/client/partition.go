package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListPartitions(metalake, catalog, schema, table string) (*models.IdentifiersResponse, error) {
	var result models.IdentifiersResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partitions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) GetPartition(metalake, catalog, schema, table, name string) (*models.PartitionResponse, error) {
	var result models.PartitionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partitions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table), url.PathEscape(name))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) CreatePartition(metalake, catalog, schema, table string, req *models.PartitionCreateRequest) (*models.PartitionResponse, error) {
	var result models.PartitionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partitions",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table))
	err := c.Post(path, req, &result)
	return &result, err
}

func (c *Client) UpdatePartition(metalake, catalog, schema, table, name string, updates []interface{}) (*models.PartitionResponse, error) {
	var result models.PartitionResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partitions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table), url.PathEscape(name))
	err := c.Put(path, &models.PartitionUpdateRequest{Updates: updates}, &result)
	return &result, err
}

func (c *Client) DropPartition(metalake, catalog, schema, table, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partitions/%s",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) ListPartitionsDetails(metalake, catalog, schema, table string) ([]models.Partition, error) {
	identifiers, err := c.ListPartitions(metalake, catalog, schema, table)
	if err != nil {
		return nil, err
	}
	partitions := make([]models.Partition, 0, len(identifiers.Identifiers))
	for _, id := range identifiers.Identifiers {
		resp, err := c.GetPartition(metalake, catalog, schema, table, id.Name)
		if err != nil {
			return nil, err
		}
		partitions = append(partitions, resp.Partition)
	}
	return partitions, nil
}
