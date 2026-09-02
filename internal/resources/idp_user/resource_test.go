package idp_user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIdpUserResource_Metadata(t *testing.T) {
	r := New()
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_idp_user" {
		t.Fatalf("expected gravitino_idp_user, got %s", resp.TypeName)
	}
}

func TestIdpUserResource_Schema(t *testing.T) {
	r := New()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	for _, name := range []string{"id", "name", "password", "enabled", "groups"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("missing attribute %s", name)
		}
	}
}
