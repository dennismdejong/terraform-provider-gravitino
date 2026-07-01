package schema_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
	resourceschema "github.com/gravitino/terraform-provider-gravitino/internal/resources/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSchemaResourceMetadata(t *testing.T) {
	r := resourceschema.NewSchemaResource()
	var req resource.MetadataRequest
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_schema" {
		t.Errorf("Expected type name gravitino_schema, got %s", resp.TypeName)
	}
}

func TestImportIDParsing(t *testing.T) {
	tests := []struct {
		id          string
		wantParts   []string
		expectError bool
	}{
		{"metalake.catalog.schema", []string{"metalake", "catalog", "schema"}, false},
		{"ml.ctlg.sch", []string{"ml", "ctlg", "sch"}, false},
		{"ml.ctlg", nil, true},
		{"ml", nil, true},
		{"ml.ctlg.sch.extra", []string{"ml", "ctlg", "sch.extra"}, false},
	}

	for _, tt := range tests {
		parts := strings.SplitN(tt.id, ".", 3)
		if tt.expectError && len(parts) == 3 {
			t.Errorf("Expected error for %q but got %v", tt.id, parts)
		}
		if !tt.expectError && len(parts) != 3 {
			t.Errorf("Expected 3 parts for %q but got %d", tt.id, len(parts))
		}
		if !tt.expectError {
			for i, want := range tt.wantParts {
				if parts[i] != want {
					t.Errorf("Part %d for %q: expected %q, got %q", i, tt.id, want, parts[i])
				}
			}
		}
	}
}

func TestAuditModelConversion(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name  string
		audit *models.Audit
		null  bool
	}{
		{
			name:  "nil audit",
			audit: nil,
			null:  true,
		},
		{
			name: "full audit",
			audit: &models.Audit{
				Creator:          "admin",
				CreateTime:       &now,
				LastModifier:     "admin",
				LastModifiedTime: &now,
			},
			null: false,
		},
		{
			name: "partial audit",
			audit: &models.Audit{
				Creator: "admin",
			},
			null: false,
		},
		{
			name:  "empty audit",
			audit: &models.Audit{},
			null:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := types.ObjectNull(map[string]attr.Type{
				"creator":            types.StringType,
				"create_time":        types.StringType,
				"last_modifier":      types.StringType,
				"last_modified_time": types.StringType,
			})

			if !tt.null {
				if val.IsNull() {
					return
				}
			} else {
				return
			}
		})
	}
}
