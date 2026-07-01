package metalake

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*MetalakeDataSource)(nil)

type MetalakeDataSource struct {
	client *client.Client
}

type MetalakeDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      *AuditModel  `tfsdk:"audit"`
}

func NewMetalakeDataSource() datasource.DataSource {
	return &MetalakeDataSource{}
}

func (d *MetalakeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metalake"
}

func (d *MetalakeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"comment": schema.StringAttribute{
				Computed: true,
			},
			"properties": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"audit": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"creator": schema.StringAttribute{
						Computed: true,
					},
					"create_time": schema.StringAttribute{
						Computed: true,
					},
					"last_modifier": schema.StringAttribute{
						Computed: true,
					},
					"last_modified_time": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *MetalakeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = cli
}

func (d *MetalakeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state MetalakeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetMetalake(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read metalake", err.Error())
		return
	}

	state.Comment = types.StringValue(result.Metalake.Comment)
	state.Properties = propertiesToMapDS(ctx, result.Metalake.Properties, &resp.Diagnostics)

	if result.Metalake.Audit != nil {
		state.Audit = &AuditModel{
			Creator:          types.StringValue(result.Metalake.Audit.Creator),
			CreateTime:       timeToStringDS(result.Metalake.Audit.CreateTime),
			LastModifier:     types.StringValue(result.Metalake.Audit.LastModifier),
			LastModifiedTime: timeToStringDS(result.Metalake.Audit.LastModifiedTime),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
