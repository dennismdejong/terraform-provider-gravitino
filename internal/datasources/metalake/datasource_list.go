package metalake

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*MetalakesDataSource)(nil)

type MetalakesDataSource struct {
	client *client.Client
}

type MetalakesDataSourceModel struct {
	Metalakes []MetalakeItemModel `tfsdk:"metalakes"`
}

type MetalakeItemModel struct {
	Name       types.String `tfsdk:"name"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      *AuditModel  `tfsdk:"audit"`
}

func NewMetalakesDataSource() datasource.DataSource {
	return &MetalakesDataSource{}
}

func (d *MetalakesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metalakes"
}

func (d *MetalakesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"metalakes": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
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
				},
			},
		},
	}
}

func (d *MetalakesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetalakesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state MetalakesDataSourceModel

	result, err := d.client.ListMetalakes()
	if err != nil {
		resp.Diagnostics.AddError("Failed to list metalakes", err.Error())
		return
	}

	for _, ml := range result.Metalakes {
		item := MetalakeItemModel{
			Name:       types.StringValue(ml.Name),
			Comment:    types.StringValue(ml.Comment),
			Properties: propertiesToMapDS(ctx, ml.Properties, &resp.Diagnostics),
		}
		if ml.Audit != nil {
			item.Audit = &AuditModel{
				Creator:          types.StringValue(ml.Audit.Creator),
				CreateTime:       timeToStringDS(ml.Audit.CreateTime),
				LastModifier:     types.StringValue(ml.Audit.LastModifier),
				LastModifiedTime: timeToStringDS(ml.Audit.LastModifiedTime),
			}
		}
		state.Metalakes = append(state.Metalakes, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
