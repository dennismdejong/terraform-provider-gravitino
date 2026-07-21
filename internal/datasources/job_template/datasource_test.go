package job_template_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/job_template"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestJobTemplatesDataSource_Schema(t *testing.T) {
	d := ds.NewJobTemplatesDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_job_templates" {
		t.Fatalf("expected gravitino_job_templates, got %s", resp.TypeName)
	}
}

func TestJobTemplatesDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/jobs/templates"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.JobTemplateListResponse{
			Code: 0,
			JobTemplates: []models.JobTemplate{
				{
					Name:       "template1",
					Template:   "CREATE TABLE ...",
					Comment:    "test template",
					Parameters: map[string]string{"key": "value"},
					Properties: map[string]string{"owner": "admin"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewJobTemplatesDataSource()
	d.(*ds.JobTemplatesDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	jobTemplateItemObjType := types.ObjectType{AttrTypes: ds.JobTemplateItemAttrTypes}
	templatesListType := types.ListType{ElemType: jobTemplateItemObjType}

	configModel := ds.JobTemplatesDataSourceModel{
		Metalake:  types.StringValue("test_metalake"),
		Templates: types.ListNull(jobTemplateItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake":  types.StringType,
		"templates": templatesListType,
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

func TestJobTemplateDataSource_Schema(t *testing.T) {
	d := ds.NewJobTemplateDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_job_template" {
		t.Fatalf("expected gravitino_job_template, got %s", resp.TypeName)
	}
}

func TestJobTemplateDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/jobs/templates/test_template"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.JobTemplateResponse{
			Code: 0,
			JobTemplate: models.JobTemplate{
				Name:       "test_template",
				Template:   "CREATE TABLE ...",
				Comment:    "a test template",
				Parameters: map[string]string{"env": "test"},
				Properties: map[string]string{"owner": "admin"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewJobTemplateDataSource()
	d.(*ds.JobTemplateDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":   types.StringType,
		"name":       types.StringType,
		"template":   types.StringType,
		"parameters": types.MapType{ElemType: types.StringType},
		"comment":    types.StringType,
		"properties": types.MapType{ElemType: types.StringType},
		"audit":      types.ObjectType{AttrTypes: ds.AuditAttrTypes},
	}

	configModel := ds.JobTemplateDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Name:       types.StringValue("test_template"),
		Audit:      types.ObjectNull(ds.AuditAttrTypes),
		Parameters: types.MapNull(types.StringType),
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
