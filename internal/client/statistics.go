package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListStatistics(metalake, resourceType, resource string) (*models.StatisticsResponse, error) {
	var result models.StatisticsResponse
	path := fmt.Sprintf("/metalakes/%s/%s/%s/statistics",
		url.PathEscape(metalake), url.PathEscape(resourceType), url.PathEscape(resource))
	err := c.Get(path, &result)
	return &result, err
}

func (c *Client) ListPartitionStatistics(metalake, catalog, schema, table string) (*models.PartitionStatisticsResponse, error) {
	var result models.PartitionStatisticsResponse
	path := fmt.Sprintf("/metalakes/%s/catalogs/%s/schemas/%s/tables/%s/partition-statistics",
		url.PathEscape(metalake), url.PathEscape(catalog), url.PathEscape(schema), url.PathEscape(table))
	err := c.Get(path, &result)
	return &result, err
}
