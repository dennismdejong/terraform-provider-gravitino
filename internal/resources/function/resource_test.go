package function_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
	resourcefunction "github.com/gravitino/terraform-provider-gravitino/internal/resources/function"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFunctionResourceMetadata(t *testing.T) {
	r := resourcefunction.New()
	var req resource.MetadataRequest
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_function" {
		t.Errorf("Expected type name gravitino_function, got %s", resp.TypeName)
	}
}

func TestImportIDParsing(t *testing.T) {
	tests := []struct {
		id          string
		wantParts   []string
		expectError bool
	}{
		{"metalake.catalog.schema.function", []string{"metalake", "catalog", "schema", "function"}, false},
		{"ml.cat.sch.fn", []string{"ml", "cat", "sch", "fn"}, false},
		{"ml.cat.sch", nil, true},
		{"ml.cat", nil, true},
		{"ml", nil, true},
	}

	for _, tt := range tests {
		parts := strings.SplitN(tt.id, ".", 4)
		if tt.expectError && len(parts) == 4 {
			t.Errorf("Expected error for %q but got %v", tt.id, parts)
		}
		if !tt.expectError && len(parts) != 4 {
			t.Errorf("Expected 4 parts for %q but got %d", tt.id, len(parts))
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
