package model_version_test

import (
	"context"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/model_version"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModelVersionDataSourceMetadata(t *testing.T) {
	d := ds.NewModelVersionDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_model_version" {
		t.Errorf("Expected type name gravitino_model_version, got %s", resp.TypeName)
	}
}

func TestModelVersionsDataSourceMetadata(t *testing.T) {
	d := ds.NewModelVersionsDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_model_versions" {
		t.Errorf("Expected type name gravitino_model_versions, got %s", resp.TypeName)
	}
}
