package model_test

import (
	"context"
	"testing"

	datasourcemodel "github.com/gravitino/terraform-provider-gravitino/internal/datasources/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModelDataSourceMetadata(t *testing.T) {
	d := datasourcemodel.NewModelDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_model" {
		t.Errorf("Expected type name gravitino_model, got %s", resp.TypeName)
	}
}

func TestModelsDataSourceMetadata(t *testing.T) {
	d := datasourcemodel.NewModelsDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_models" {
		t.Errorf("Expected type name gravitino_models, got %s", resp.TypeName)
	}
}
