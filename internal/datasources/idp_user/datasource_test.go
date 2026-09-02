package idp_user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestIdpUserDataSource_Metadata(t *testing.T) {
	d := NewDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_idp_user" {
		t.Fatalf("expected gravitino_idp_user, got %s", resp.TypeName)
	}
}

func TestIdpUserDataSource_Schema(t *testing.T) {
	d := NewDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	for _, name := range []string{"name", "enabled", "groups"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("missing attribute %s", name)
		}
	}
}
