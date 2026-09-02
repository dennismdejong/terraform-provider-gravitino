package idp_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIdpGroupResource_Metadata(t *testing.T) {
	r := New()
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_idp_group" {
		t.Fatalf("expected gravitino_idp_group, got %s", resp.TypeName)
	}
}

func TestIdpGroupResource_Schema(t *testing.T) {
	r := New()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	for _, name := range []string{"id", "name", "comment", "users"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("missing attribute %s", name)
		}
	}
}
