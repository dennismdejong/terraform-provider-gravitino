package catalog_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/catalog"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCatalogsDataSource_Schema(t *testing.T) {
	d := ds.NewListDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_catalogs" {
		t.Fatalf("expected gravitino_catalogs, got %s", resp.TypeName)
	}
}

func TestCatalogsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("details") != "true" {
			t.Errorf("expected details=true, got %s", r.URL.Query().Get("details"))
		}

		resp := models.CatalogInfoListResponse{
			Code: 0,
			Catalogs: []models.Catalog{
				{
					Name:       "catalog1",
					Type:       "relational",
					Provider:   "hive",
					Comment:    "test catalog",
					Properties: map[string]string{"key": "value"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewListDataSource()
	d.(*ds.CatalogsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	catItemObjType := types.ObjectType{AttrTypes: ds.CatalogItemAttrTypes}
	catalogsListType := types.ListType{ElemType: catItemObjType}

	configModel := ds.CatalogsDataSourceModel{
		Metalake: types.StringValue("test_metalake"),
		Catalogs: types.ListNull(catItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake": types.StringType,
		"catalogs": catalogsListType,
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

func TestCatalogDataSource_Schema(t *testing.T) {
	d := ds.NewGetDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_catalog" {
		t.Fatalf("expected gravitino_catalog, got %s", resp.TypeName)
	}
}

func TestCatalogDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.CatalogResponse{
			Code: 0,
			Catalog: models.Catalog{
				Name:       "test_catalog",
				Type:       "relational",
				Provider:   "hive",
				Comment:    "a test catalog",
				Properties: map[string]string{"env": "test"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewGetDataSource()
	d.(*ds.CatalogDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":   types.StringType,
		"name":       types.StringType,
		"type":       types.StringType,
		"catalog_provider":   types.StringType,
		"comment":    types.StringType,
		"properties": types.MapType{ElemType: types.StringType},
		"audit":      types.ObjectType{AttrTypes: ds.AuditAttrTypes},
	}

	configModel := ds.CatalogDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Name:       types.StringValue("test_catalog"),
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
