package model_test

import (
	"context"
	"testing"

	resourcemodel "github.com/gravitino/terraform-provider-gravitino/internal/resources/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestModelResourceMetadata(t *testing.T) {
	r := resourcemodel.New()
	var req resource.MetadataRequest
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_model" {
		t.Errorf("Expected type name gravitino_model, got %s", resp.TypeName)
	}
}
