package policy_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/policy"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPoliciesDataSource_Schema(t *testing.T) {
	d := ds.NewListDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_policies" {
		t.Fatalf("expected gravitino_policies, got %s", resp.TypeName)
	}
}

func TestPoliciesDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.PolicyListResponse{
			Code: 0,
			Policies: []models.Policy{
				{
					Name:       "policy1",
					Condition:  "",
					Effect:     "allow",
					Actions:    []string{"read", "write"},
					Subjects:   []string{"user1"},
					Object:     "",
					Properties: map[string]string{"key": "value"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewListDataSource()
	d.(*ds.PoliciesDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	policyItemObjType := types.ObjectType{AttrTypes: ds.PolicyItemAttrTypes}
	policiesListType := types.ListType{ElemType: policyItemObjType}

	configModel := ds.PoliciesDataSourceModel{
		Metalake:     types.StringValue("test_metalake"),
		ResourceType: types.StringValue("catalogs"),
		Resource:     types.StringValue("test_catalog"),
		Policies:     types.ListNull(policyItemObjType),
	}

	attrTypes := map[string]attr.Type{
		"metalake":      types.StringType,
		"resource_type": types.StringType,
		"resource":      types.StringType,
		"policies":      policiesListType,
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
