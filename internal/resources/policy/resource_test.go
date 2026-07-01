package policy_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
	res "github.com/gravitino/terraform-provider-gravitino/internal/resources/policy"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPolicyResource_Schema(t *testing.T) {
	r := res.New()
	resp := &resource.MetadataResponse{}
	r.Metadata(nil, resource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_policy" {
		t.Fatalf("expected gravitino_policy, got %s", resp.TypeName)
	}
}

func TestPolicyResource_Create(t *testing.T) {
	var receivedBody models.PolicyCreateRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)

		resp := models.PolicyResponse{
			Code: 0,
			Policy: models.Policy{
				Name:       receivedBody.Name,
				Condition:  receivedBody.Condition,
				Effect:     receivedBody.Effect,
				Actions:    receivedBody.Actions,
				Subjects:   receivedBody.Subjects,
				Object:     receivedBody.Object,
				Properties: receivedBody.Properties,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	r := res.New()
	r.(*res.PolicyResource).SetClient(c)

	ctx := context.Background()

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	planModel := res.PolicyResourceModel{
		Metalake:     types.StringValue("test_metalake"),
		ResourceType: types.StringValue("catalogs"),
		Resource:     types.StringValue("test_catalog"),
		Name:         types.StringValue("test_policy"),
		Effect:       types.StringValue("allow"),
		Condition:    types.StringNull(),
		Actions:      types.ListNull(types.StringType),
		Subjects:     types.ListNull(types.StringType),
		Object:       types.StringNull(),
		Properties:   types.MapNull(types.StringType),
		Audit:        types.ObjectNull(res.AuditAttrTypes),
	}

	planObj, diags := types.ObjectValueFrom(ctx, schemaObj.Type().(types.ObjectType).AttributeTypes(), planModel)
	if diags.HasError() {
		t.Fatalf("failed to create plan object: %v", diags)
	}

	tfVal, err := planObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := resource.CreateRequest{
		Plan: tfsdk.Plan{Schema: schemaObj, Raw: tfVal},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{Schema: schemaObj},
	}

	r.Create(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			t.Logf("diag error: %s: %s", diag.Summary(), diag.Detail())
		}
		t.Fatal("unexpected diagnostics errors")
	}

	if receivedBody.Name != "test_policy" {
		t.Errorf("expected name test_policy, got %s", receivedBody.Name)
	}
}

func TestPolicyResource_ImportState(t *testing.T) {
	r := res.New().(resource.ResourceWithImportState)

	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	r.(resource.Resource).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	nullModel := res.PolicyResourceModel{
		Metalake:     types.StringNull(),
		ResourceType: types.StringNull(),
		Resource:     types.StringNull(),
		Name:         types.StringNull(),
		Effect:       types.StringNull(),
		Condition:    types.StringNull(),
		Object:       types.StringNull(),
		Actions:      types.ListNull(types.StringType),
		Subjects:     types.ListNull(types.StringType),
		Properties:   types.MapNull(types.StringType),
		Audit:        types.ObjectNull(res.AuditAttrTypes),
	}
	nullObj, diags := types.ObjectValueFrom(ctx, schemaObj.Type().(types.ObjectType).AttributeTypes(), nullModel)
	if diags.HasError() {
		t.Fatalf("failed to create null object: %v", diags)
	}
	rawVal, err := nullObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := resource.ImportStateRequest{
		ID: "my_metalake.catalogs.my_catalog.my_policy",
	}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{Schema: schemaObj, Raw: rawVal},
	}

	r.ImportState(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			t.Logf("diag error: %s: %s", diag.Summary(), diag.Detail())
		}
		t.Fatal("unexpected diagnostics errors")
	}
}

func TestPolicyResource_ImportState_Invalid(t *testing.T) {
	r := res.New().(resource.ResourceWithImportState)

	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	r.(resource.Resource).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	nullModel := res.PolicyResourceModel{
		Metalake:     types.StringNull(),
		ResourceType: types.StringNull(),
		Resource:     types.StringNull(),
		Name:         types.StringNull(),
		Effect:       types.StringNull(),
		Condition:    types.StringNull(),
		Object:       types.StringNull(),
		Actions:      types.ListNull(types.StringType),
		Subjects:     types.ListNull(types.StringType),
		Properties:   types.MapNull(types.StringType),
		Audit:        types.ObjectNull(res.AuditAttrTypes),
	}
	nullObj, diags := types.ObjectValueFrom(ctx, schemaObj.Type().(types.ObjectType).AttributeTypes(), nullModel)
	if diags.HasError() {
		t.Fatalf("failed to create null object: %v", diags)
	}
	rawVal, err := nullObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := resource.ImportStateRequest{
		ID: "too.few.dots",
	}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{Schema: schemaObj, Raw: rawVal},
	}

	r.ImportState(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for invalid import ID")
	}
}

func TestPolicyResource_Delete(t *testing.T) {
	var deleteCalled bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteCalled = true
		}
		resp := models.DropResponse{Code: 0, Dropped: true}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	r := res.New()
	r.(*res.PolicyResource).SetClient(c)

	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	stateModel := res.PolicyResourceModel{
		ID:           types.StringValue("test_metalake.catalogs.test_catalog.test_policy"),
		Metalake:     types.StringValue("test_metalake"),
		ResourceType: types.StringValue("catalogs"),
		Resource:     types.StringValue("test_catalog"),
		Name:         types.StringValue("test_policy"),
		Effect:       types.StringValue("allow"),
		Condition:    types.StringNull(),
		Object:       types.StringNull(),
		Actions:      types.ListNull(types.StringType),
		Subjects:     types.ListNull(types.StringType),
		Properties:   types.MapNull(types.StringType),
		Audit:        types.ObjectNull(res.AuditAttrTypes),
	}

	stateObj, diags := types.ObjectValueFrom(ctx, schemaObj.Type().(types.ObjectType).AttributeTypes(), stateModel)
	if diags.HasError() {
		t.Fatalf("failed to create state object: %v", diags)
	}

	tfVal, err := stateObj.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	req := resource.DeleteRequest{
		State: tfsdk.State{Schema: schemaObj, Raw: tfVal},
	}
	resp := &resource.DeleteResponse{
		State: tfsdk.State{Schema: schemaObj},
	}

	r.Delete(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			t.Logf("diag error: %s: %s", diag.Summary(), diag.Detail())
		}
		t.Fatal("unexpected diagnostics errors")
	}

	if !deleteCalled {
		t.Fatal("delete was not called")
	}
}
