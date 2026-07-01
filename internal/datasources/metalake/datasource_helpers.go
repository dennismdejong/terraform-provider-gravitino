package metalake

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AuditModel struct {
	Creator          types.String `tfsdk:"creator"`
	CreateTime       types.String `tfsdk:"create_time"`
	LastModifier     types.String `tfsdk:"last_modifier"`
	LastModifiedTime types.String `tfsdk:"last_modified_time"`
}

func propertiesToMapDS(ctx context.Context, props map[string]string, diags *diag.Diagnostics) types.Map {
	if len(props) == 0 {
		return types.MapNull(types.StringType)
	}
	result, d := types.MapValueFrom(ctx, types.StringType, props)
	if d.HasError() {
		*diags = append(*diags, d...)
		return types.MapNull(types.StringType)
	}
	return result
}

func timeToStringDS(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}
