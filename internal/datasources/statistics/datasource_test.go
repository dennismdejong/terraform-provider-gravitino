package statistics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/statistics"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStatisticsDataSource_Schema(t *testing.T) {
	d := ds.New()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_statistics" {
		t.Fatalf("expected gravitino_statistics, got %s", resp.TypeName)
	}
}

func TestStatisticsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog/statistics"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.StatisticsResponse{
			Code: 0,
			Statistics: []models.Statistics{
				{
					Name:       "numRows",
					Type:       "long",
					Value:      "1000",
					Properties: map[string]string{"unit": "rows"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.New()
	d.(*ds.StatisticsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	statItemObjType := types.ObjectType{AttrTypes: ds.StatisticItemAttrTypes}
	statsListType := types.ListType{ElemType: statItemObjType}

	attrTypes := map[string]attr.Type{
		"metalake":      types.StringType,
		"resource_type": types.StringType,
		"resource":      types.StringType,
		"statistics":    statsListType,
	}

	configModel := ds.StatisticsDataSourceModel{
		Metalake:     types.StringValue("test_metalake"),
		ResourceType: types.StringValue("catalogs"),
		Resource:     types.StringValue("test_catalog"),
		Statistics:   types.ListNull(statItemObjType),
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

func TestPartitionStatisticsDataSource_Schema(t *testing.T) {
	d := ds.NewPartition()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_partition_statistics" {
		t.Fatalf("expected gravitino_partition_statistics, got %s", resp.TypeName)
	}
}

func TestPartitionStatisticsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog/schemas/test_schema/tables/test_table/partition-statistics"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.PartitionStatisticsResponse{
			Code: 0,
			Statistics: []models.PartitionStatistics{
				{
					PartitionName: "dt=2024-01-01",
					Statistics: []models.Statistics{
						{
							Name:       "numRows",
							Type:       "long",
							Value:      "500",
							Properties: map[string]string{"unit": "rows"},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewPartition()
	d.(*ds.PartitionStatisticsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	psItemObjType := types.ObjectType{AttrTypes: ds.PartitionStatisticItemAttrTypes}
	psListType := types.ListType{ElemType: psItemObjType}

	attrTypes := map[string]attr.Type{
		"metalake":             types.StringType,
		"catalog":              types.StringType,
		"schema":               types.StringType,
		"table":                types.StringType,
		"partition_statistics": psListType,
	}

	configModel := ds.PartitionStatisticsDataSourceModel{
		Metalake:            types.StringValue("test_metalake"),
		Catalog:             types.StringValue("test_catalog"),
		Schema:              types.StringValue("test_schema"),
		Table:               types.StringValue("test_table"),
		PartitionStatistics: types.ListNull(psItemObjType),
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
