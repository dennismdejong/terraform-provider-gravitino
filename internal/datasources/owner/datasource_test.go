package owner_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/owner"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOwnerDataSource_Schema(t *testing.T) {
	d := ds.NewOwnerDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_owner" {
		t.Fatalf("expected gravitino_owner, got %s", resp.TypeName)
	}
}

func TestOwnerDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/owners/CATALOG/test_catalog"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.OwnerResponse{
			Code: 0,
			Owner: models.Owner{
				Name: "admin",
				Type: "USER",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, "", "", "", "")
	d := ds.NewOwnerDataSource()
	d.(*ds.OwnerDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":         types.StringType,
		"object_type":      types.StringType,
		"object_full_name": types.StringType,
		"owner_name":       types.StringType,
		"owner_type":       types.StringType,
	}

	configModel := ds.OwnerDataSourceModel{
		Metalake:       types.StringValue("test_metalake"),
		ObjectType:     types.StringValue("CATALOG"),
		ObjectFullName: types.StringValue("test_catalog"),
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
