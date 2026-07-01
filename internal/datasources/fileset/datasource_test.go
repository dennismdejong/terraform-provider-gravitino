package fileset_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/fileset"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFilesetsDataSource_Schema(t *testing.T) {
	d := ds.NewListDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_filesets" {
		t.Fatalf("expected gravitino_filesets, got %s", resp.TypeName)
	}
}

func TestFilesetsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/filesets" {
			resp := models.IdentifiersResponse{
				Code: 0,
				Identifiers: []models.NameIdentifier{
					{Namespace: []string{"test_metalake", "test_catalog", "test_schema"}, Name: "fileset1"},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/filesets/fileset1" {
			resp := models.FilesetResponse{
				Code: 0,
				Fileset: models.Fileset{
					Name:       "fileset1",
					Type:       "managed",
					Comment:    "test fileset",
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
	d := ds.NewListDataSource()
	d.(*ds.FilesetsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	fsItemObjType := types.ObjectType{AttrTypes: ds.FilesetItemAttrTypes}
	filesetsListType := types.ListType{ElemType: fsItemObjType}

	configModel := ds.FilesetsDataSourceModel{
		Metalake: types.StringValue("test_metalake"),
		Catalog:  types.StringValue("test_catalog"),
		Schema:   types.StringValue("test_schema"),
		Filesets: types.ListNull(fsItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake": types.StringType,
		"catalog":  types.StringType,
		"schema":   types.StringType,
		"filesets": filesetsListType,
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

func TestFilesetDataSource_Schema(t *testing.T) {
	d := ds.NewGetDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_fileset" {
		t.Fatalf("expected gravitino_fileset, got %s", resp.TypeName)
	}
}

func TestFilesetDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/filesets/test_fileset"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.FilesetResponse{
			Code: 0,
			Fileset: models.Fileset{
				Name:            "test_fileset",
				Type:            "managed",
				Comment:         "a test fileset",
				StorageLocation: "s3://bucket/path",
				Properties:      map[string]string{"env": "test"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewGetDataSource()
	d.(*ds.FilesetDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":         types.StringType,
		"catalog":          types.StringType,
		"schema":           types.StringType,
		"name":             types.StringType,
		"comment":          types.StringType,
		"type":             types.StringType,
		"storage_location": types.StringType,
		"properties":       types.MapType{ElemType: types.StringType},
		"audit":            types.ObjectType{AttrTypes: ds.AuditAttrTypes},
	}

	configModel := ds.FilesetDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Catalog:    types.StringValue("test_catalog"),
		Schema:     types.StringValue("test_schema"),
		Name:       types.StringValue("test_fileset"),
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
