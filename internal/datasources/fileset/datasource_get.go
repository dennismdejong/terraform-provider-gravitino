package fileset

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

var _ datasource.DataSource = &FilesetDataSource{}
var _ datasource.DataSourceWithConfigure = &FilesetDataSource{}

type FilesetDataSource struct {
	client *client.Client
}

func NewGetDataSource() datasource.DataSource {
	return &FilesetDataSource{}
}

func (d *FilesetDataSource) SetClient(c *client.Client) {
	d.client = c
}

type FilesetDataSourceModel struct {
	Metalake        types.String `tfsdk:"metalake"`
	Catalog         types.String `tfsdk:"catalog"`
	Schema          types.String `tfsdk:"schema"`
	Name            types.String `tfsdk:"name"`
	Comment         types.String `tfsdk:"comment"`
	Type            types.String `tfsdk:"type"`
	StorageLocation types.String `tfsdk:"storage_location"`
	Properties      types.Map    `tfsdk:"properties"`
	Audit           types.Object `tfsdk:"audit"`
}

func (d *FilesetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FilesetDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_fileset"
}

func (d *FilesetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"catalog": schema.StringAttribute{
				Required:    true,
				Description: "The catalog name.",
			},
			"schema": schema.StringAttribute{
				Required:    true,
				Description: "The schema name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The fileset name.",
			},
			"comment": schema.StringAttribute{
				Computed:    true,
				Description: "The fileset comment.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The fileset type.",
			},
			"storage_location": schema.StringAttribute{
				Computed:    true,
				Description: "The fileset storage location.",
			},
			"properties": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The fileset properties.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the fileset.",
			},
		},
	}
}

func (d *FilesetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FilesetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetFileset(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get fileset", err.Error())
		return
	}

	setDataSourceStateFromFileset(ctx, &resp.Diagnostics, &result.Fileset, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func setDataSourceStateFromFileset(_ context.Context, diags *diag.Diagnostics, fileset *models.Fileset, model *FilesetDataSourceModel) {
	model.Comment = types.StringValue(fileset.Comment)
	model.Type = types.StringValue(fileset.Type)
	model.StorageLocation = types.StringValue(fileset.StorageLocation)

	props, d := types.MapValueFrom(nil, types.StringType, fileset.Properties)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	model.Properties = props

	model.Audit = auditToObjectValueForDS(nil, fileset.Audit)
}
