package partition_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/partition"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPartitionsDataSource_Schema(t *testing.T) {
	d := ds.NewPartitionsDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_partitions" {
		t.Fatalf("expected gravitino_partitions, got %s", resp.TypeName)
	}
}

func TestPartitionsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/tables/test_table/partitions" {
			resp := models.IdentifiersResponse{
				Code: 0,
				Identifiers: []models.NameIdentifier{
					{Name: "p1", Namespace: []string{"test_metalake", "test_catalog", "test_schema", "test_table"}},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/tables/test_table/partitions/p1" {
			resp := models.PartitionResponse{
				Code: 0,
				Partition: models.Partition{
					Name:       "p1",
					Properties: map[string]string{"key": "value"},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewPartitionsDataSource().(*ds.PartitionsDataSource)
	d.SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	catItemObjType := types.ObjectType{AttrTypes: ds.PartitionItemAttrTypes}
	partitionsListType := types.ListType{ElemType: catItemObjType}

	configModel := ds.PartitionsDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Catalog:    types.StringValue("test_catalog"),
		Schema:     types.StringValue("test_schema"),
		Table:      types.StringValue("test_table"),
		Partitions: types.ListNull(catItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake":   types.StringType,
		"catalog":    types.StringType,
		"schema":     types.StringType,
		"table":      types.StringType,
		"partitions": partitionsListType,
	}

	configObj, diags := types.ObjectValueFrom(ctx, attrTypes, configModel)
	if diags.HasError() {
		t.Fatalf("failed to create config object: %v", diags)
	}

	tfVal, err := configObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := datasource.ReadRequest{
		Config: tfsdk.Config{Schema: schemaObj, Raw: tfVal},
	}
	resp := &datasource.ReadResponse{
		State: tfsdk.State{Schema: schemaObj},
	}

	d.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			t.Logf("diag error: %s: %s", diag.Summary(), diag.Detail())
		}
		t.Fatal("unexpected diagnostics errors")
	}
}

func TestPartitionDataSource_Schema(t *testing.T) {
	d := ds.NewPartitionDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_partition" {
		t.Fatalf("expected gravitino_partition, got %s", resp.TypeName)
	}
}

func TestPartitionDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/tables/test_table/partitions/test_partition"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.PartitionResponse{
			Code: 0,
			Partition: models.Partition{
				Name:       "test_partition",
				Properties: map[string]string{"env": "test"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewPartitionDataSource().(*ds.PartitionDataSource)
	d.SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":   types.StringType,
		"catalog":    types.StringType,
		"schema":     types.StringType,
		"table":      types.StringType,
		"name":       types.StringType,
		"properties": types.MapType{ElemType: types.StringType},
		"audit":      types.ObjectType{AttrTypes: ds.AuditAttrTypes},
	}

	configModel := ds.PartitionDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Catalog:    types.StringValue("test_catalog"),
		Schema:     types.StringValue("test_schema"),
		Table:      types.StringValue("test_table"),
		Name:       types.StringValue("test_partition"),
		Audit:      types.ObjectNull(ds.AuditAttrTypes),
		Properties: types.MapNull(types.StringType),
	}

	configObj, diags := types.ObjectValueFrom(ctx, attrTypes, configModel)
	if diags.HasError() {
		t.Fatalf("failed to create config object: %v", diags)
	}

	tfVal, err := configObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := datasource.ReadRequest{
		Config: tfsdk.Config{Schema: schemaObj, Raw: tfVal},
	}
	resp := &datasource.ReadResponse{
		State: tfsdk.State{Schema: schemaObj},
	}

	d.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			t.Logf("diag error: %s: %s", diag.Summary(), diag.Detail())
		}
		t.Fatal("unexpected diagnostics errors")
	}
}
