package job

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &JobDataSource{}
var _ datasource.DataSourceWithConfigure = &JobDataSource{}

type JobDataSource struct {
	client *client.Client
}

func NewGetDataSource() datasource.DataSource {
	return &JobDataSource{}
}

func (d *JobDataSource) SetClient(c *client.Client) {
	d.client = c
}

type JobDataSourceModel struct {
	Metalake   types.String `tfsdk:"metalake"`
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Parameters types.Map    `tfsdk:"parameters"`
	Schedule   types.String `tfsdk:"schedule"`
	Status     types.String `tfsdk:"status"`
	Audit      types.Object `tfsdk:"audit"`
}

func (d *JobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *JobDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_job"
}

func (d *JobDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The job name.",
			},
			"template": schema.StringAttribute{
				Computed:    true,
				Description: "The job template name.",
			},
			"parameters": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The job parameters.",
			},
			"schedule": schema.StringAttribute{
				Computed:    true,
				Description: "The job schedule.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the job.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the job.",
			},
		},
	}
}

func (d *JobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JobDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetJob(config.Metalake.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get job", err.Error())
		return
	}

	setDataSourceStateFromJob(ctx, &resp.Diagnostics, &result.Job, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func setDataSourceStateFromJob(ctx context.Context, diags *diag.Diagnostics, job *models.Job, model *JobDataSourceModel) {
	model.Template = types.StringValue(job.Template)
	model.Schedule = types.StringValue(job.Schedule)
	model.Status = types.StringValue(job.Status)

	props, d := types.MapValueFrom(ctx, types.StringType, strMapFromInterface(job.Parameters))
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	model.Parameters = props

	model.Audit = auditToObjectValueForDS(ctx, job.Audit)
}
