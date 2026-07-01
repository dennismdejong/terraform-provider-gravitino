package view_test

import (
	"context"
	"testing"

	datasourceview "github.com/gravitino/terraform-provider-gravitino/internal/datasources/view"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestViewDataSourceMetadata(t *testing.T) {
	d := datasourceview.NewViewDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_view" {
		t.Errorf("Expected type name gravitino_view, got %s", resp.TypeName)
	}
}

func TestViewsDataSourceMetadata(t *testing.T) {
	d := datasourceview.NewViewsDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_views" {
		t.Errorf("Expected type name gravitino_views, got %s", resp.TypeName)
	}
}
