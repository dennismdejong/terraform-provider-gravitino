package function_test

import (
	"context"
	"testing"

	datasourcefunction "github.com/gravitino/terraform-provider-gravitino/internal/datasources/function"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestFunctionDataSourceMetadata(t *testing.T) {
	d := datasourcefunction.NewFunctionDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_function" {
		t.Errorf("Expected type name gravitino_function, got %s", resp.TypeName)
	}
}

func TestFunctionsDataSourceMetadata(t *testing.T) {
	d := datasourcefunction.NewFunctionsDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_functions" {
		t.Errorf("Expected type name gravitino_functions, got %s", resp.TypeName)
	}
}
