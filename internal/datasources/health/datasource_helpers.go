package health

import (
	"context"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func healthChecksToModel(ctx context.Context, checks []models.HealthCheck, diags *diag.Diagnostics) []HealthCheckModel {
	if checks == nil {
		return nil
	}
	result := make([]HealthCheckModel, len(checks))
	for i, c := range checks {
		details, d := types.MapValueFrom(ctx, types.StringType, c.Details)
		diags.Append(d...)
		if diags.HasError() {
			return nil
		}
		result[i] = HealthCheckModel{
			Name:    types.StringValue(c.Name),
			Status:  types.StringValue(c.Status),
			Details: details,
		}
	}
	return result
}
