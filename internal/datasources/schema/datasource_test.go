package schema_test

import (
	"context"
	"testing"

	datasourceschema "github.com/gravitino/terraform-provider-gravitino/internal/datasources/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestSchemaDataSourceMetadata(t *testing.T) {
	d := datasourceschema.NewSchemaDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_schema" {
		t.Errorf("Expected type name gravitino_schema, got %s", resp.TypeName)
	}
}

func TestSchemasDataSourceMetadata(t *testing.T) {
	d := datasourceschema.NewSchemasDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_schemas" {
		t.Errorf("Expected type name gravitino_schemas, got %s", resp.TypeName)
	}
}
