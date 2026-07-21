package job_template

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

var _ datasource.DataSource = &JobTemplateDataSource{}
var _ datasource.DataSourceWithConfigure = &JobTemplateDataSource{}

type JobTemplateDataSource struct {
	client *client.Client
}

func NewJobTemplateDataSource() datasource.DataSource {
	return &JobTemplateDataSource{}
}

func (d *JobTemplateDataSource) SetClient(c *client.Client) {
	d.client = c
}

type JobTemplateDataSourceModel struct {
	Metalake   types.String `tfsdk:"metalake"`
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Parameters types.Map    `tfsdk:"parameters"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      types.Object `tfsdk:"audit"`
}

func (d *JobTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JobTemplateDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_job_template"
}

func (d *JobTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The job template name.",
			},
			"template": schema.StringAttribute{
				Computed:    true,
				Description: "The template definition.",
			},
			"parameters": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The job template parameters.",
			},
			"comment": schema.StringAttribute{
				Computed:    true,
				Description: "The job template comment.",
			},
			"properties": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The job template properties.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the job template.",
			},
		},
	}
}

func (d *JobTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JobTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetJobTemplate(config.Metalake.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get job template", err.Error())
		return
	}

	setDataSourceStateFromJobTemplate(ctx, &resp.Diagnostics, &result.JobTemplate, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func setDataSourceStateFromJobTemplate(_ context.Context, diags *diag.Diagnostics, jt *models.JobTemplate, model *JobTemplateDataSourceModel) {
	model.Template = types.StringValue(jt.Template)
	model.Comment = types.StringValue(jt.Comment)

	params, d := types.MapValueFrom(context.TODO(), types.StringType, jt.Parameters)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	model.Parameters = params

	props, d := types.MapValueFrom(context.TODO(), types.StringType, jt.Properties)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	model.Properties = props

	auditObj, d := auditToObjectValueForDS(context.TODO(), jt.Audit)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	model.Audit = auditObj
}
