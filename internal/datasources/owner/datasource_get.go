package owner

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &OwnerDataSource{}
var _ datasource.DataSourceWithConfigure = &OwnerDataSource{}

type OwnerDataSource struct {
	client *client.Client
}

func NewOwnerDataSource() datasource.DataSource {
	return &OwnerDataSource{}
}

func (d *OwnerDataSource) SetClient(c *client.Client) {
	d.client = c
}

type OwnerDataSourceModel struct {
	Metalake       types.String `tfsdk:"metalake"`
	ObjectType     types.String `tfsdk:"object_type"`
	ObjectFullName types.String `tfsdk:"object_full_name"`
	OwnerName      types.String `tfsdk:"owner_name"`
	OwnerType      types.String `tfsdk:"owner_type"`
}

func (d *OwnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OwnerDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_owner"
}

func (d *OwnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"object_type": schema.StringAttribute{
				Required:    true,
				Description: "The object type (e.g. CATALOG, SCHEMA, TABLE, etc.).",
			},
			"object_full_name": schema.StringAttribute{
				Required:    true,
				Description: "The full object name (dot-separated).",
			},
			"owner_name": schema.StringAttribute{
				Computed:    true,
				Description: "The owner name.",
			},
			"owner_type": schema.StringAttribute{
				Computed:    true,
				Description: "The owner type (USER or GROUP).",
			},
		},
	}
}

func (d *OwnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OwnerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetOwner(
		config.Metalake.ValueString(),
		config.ObjectType.ValueString(),
		config.ObjectFullName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get owner", err.Error())
		return
	}

	config.OwnerName = types.StringValue(result.Owner.Name)
	config.OwnerType = types.StringValue(result.Owner.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
