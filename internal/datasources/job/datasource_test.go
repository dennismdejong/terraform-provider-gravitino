package job_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/job"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestJobsDataSource_Schema(t *testing.T) {
	d := ds.NewListDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_jobs" {
		t.Fatalf("expected gravitino_jobs, got %s", resp.TypeName)
	}
}

func TestJobsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/jobs"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.JobListResponse{
			Code: 0,
			Jobs: []models.Job{
				{
					Name:     "job1",
					Template: "ingestion",
					Schedule: "0 0 * * *",
					Status:   "RUNNING",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewListDataSource()
	d.(*ds.JobsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	jobItemObjType := types.ObjectType{AttrTypes: ds.JobItemAttrTypes}
	jobsListType := types.ListType{ElemType: jobItemObjType}

	configModel := ds.JobsDataSourceModel{
		Metalake: types.StringValue("test_metalake"),
		Jobs:     types.ListNull(jobItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake": types.StringType,
		"jobs":     jobsListType,
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

func TestJobDataSource_Schema(t *testing.T) {
	d := ds.NewGetDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_job" {
		t.Fatalf("expected gravitino_job, got %s", resp.TypeName)
	}
}

func TestJobDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/jobs/test_job"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.JobResponse{
			Code: 0,
			Job: models.Job{
				Name:       "test_job",
				Template:   "ingestion",
				Schedule:   "0 0 * * *",
				Status:     "RUNNING",
				Parameters: map[string]interface{}{"key": "value"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewGetDataSource()
	d.(*ds.JobDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":   types.StringType,
		"name":       types.StringType,
		"template":   types.StringType,
		"parameters": types.MapType{ElemType: types.StringType},
		"schedule":   types.StringType,
		"status":     types.StringType,
		"audit":      types.ObjectType{AttrTypes: ds.AuditAttrTypes},
	}

	configModel := ds.JobDataSourceModel{
		Metalake:   types.StringValue("test_metalake"),
		Name:       types.StringValue("test_job"),
		Audit:      types.ObjectNull(ds.AuditAttrTypes),
		Parameters: types.MapNull(types.StringType),
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
