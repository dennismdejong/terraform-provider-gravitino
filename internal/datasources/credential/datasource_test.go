package credential_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "github.com/gravitino/terraform-provider-gravitino/internal/datasources/credential"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCredentialsDataSource_Schema(t *testing.T) {
	d := ds.New()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.TODO(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_credentials" {
		t.Fatalf("expected gravitino_credentials, got %s", resp.TypeName)
	}
}

func TestCredentialsDataSource_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/metalakes/test_metalake/catalogs/test_catalog/credentials"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := models.CredentialResponse{
			Code: 0,
			Credential: models.Credential{
				Type:       "s3-token",
				Value:      "my-secret-token",
				ExpireTime: "2025-12-31T23:59:59Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c, _ := client.New(server.URL, nil)
	d := ds.New()
	d.(*ds.CredentialsDataSource).SetClient(c)

	ctx := context.Background()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	schemaObj := schemaResp.Schema

	attrTypes := map[string]attr.Type{
		"metalake":      types.StringType,
		"resource_type": types.StringType,
		"resource":      types.StringType,
		"type":          types.StringType,
		"value":         types.StringType,
		"expire_time":   types.StringType,
	}

	configModel := ds.CredentialsDataSourceModel{
		Metalake:     types.StringValue("test_metalake"),
		ResourceType: types.StringValue("catalogs"),
		Resource:     types.StringValue("test_catalog"),
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
