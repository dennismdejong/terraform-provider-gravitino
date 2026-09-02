package idp_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestIdpGroupDataSource_Metadata(t *testing.T) {
	d := NewDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "gravitino_idp_group" {
		t.Fatalf("expected gravitino_idp_group, got %s", resp.TypeName)
	}
}

func TestIdpGroupDataSource_Schema(t *testing.T) {
	d := NewDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	for _, name := range []string{"name", "comment", "users"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("missing attribute %s", name)
		}
	}
}
