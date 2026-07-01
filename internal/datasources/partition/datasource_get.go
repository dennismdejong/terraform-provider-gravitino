package partition

import (
	"context"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ datasource.DataSource = &PartitionDataSource{}
var _ datasource.DataSourceWithConfigure = &PartitionDataSource{}

type PartitionDataSource struct {
	client *client.Client
}

type PartitionDataSourceModel struct {
	Metalake   types.String `tfsdk:"metalake"`
	Catalog    types.String `tfsdk:"catalog"`
	Schema     types.String `tfsdk:"schema"`
	Table      types.String `tfsdk:"table"`
	Name       types.String `tfsdk:"name"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      types.Object `tfsdk:"audit"`
}

func NewPartitionDataSource() datasource.DataSource {
	return &PartitionDataSource{}
}

func (d *PartitionDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_partition"
}

func (d *PartitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single Gravitino partition by name.",
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Description: "The metalake name.",
				Required:    true,
			},
			"catalog": schema.StringAttribute{
				Description: "The catalog name.",
				Required:    true,
			},
			"schema": schema.StringAttribute{
				Description: "The schema name.",
				Required:    true,
			},
			"table": schema.StringAttribute{
				Description: "The table name.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The partition name.",
				Required:    true,
			},
			"properties": schema.MapAttribute{
				Description: "Key-value properties for the partition.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"audit": schema.ObjectAttribute{
				Description:    "Audit information for the partition.",
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
			},
		},
	}
}

func (d *PartitionDataSource) SetClient(c *client.Client) {
	d.client = c
}

func (d *PartitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Invalid provider data", "Expected *client.Client, got unexpected type.")
		return
	}
	d.client = c
}

func (d *PartitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PartitionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	partitionResp, err := d.client.GetPartition(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString(), config.Table.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read partition", err.Error())
		return
	}

	if len(partitionResp.Partition.Properties) > 0 {
		props, d := types.MapValueFrom(ctx, types.StringType, partitionResp.Partition.Properties)
		resp.Diagnostics.Append(d...)
		config.Properties = props
	} else {
		config.Properties = types.MapNull(types.StringType)
	}

	auditObj, auditDiags := dsAuditToObject(partitionResp.Partition.Audit)
	resp.Diagnostics.Append(auditDiags...)
	config.Audit = auditObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func dsAuditToObject(audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(AuditAttrTypes), nil
	}

	creator := types.StringNull()
	if audit.Creator != "" {
		creator = types.StringValue(audit.Creator)
	}

	createTime := types.StringNull()
	if audit.CreateTime != nil {
		createTime = types.StringValue(audit.CreateTime.Format(time.RFC3339))
	}

	lastModifier := types.StringNull()
	if audit.LastModifier != "" {
		lastModifier = types.StringValue(audit.LastModifier)
	}

	lastModifiedTime := types.StringNull()
	if audit.LastModifiedTime != nil {
		lastModifiedTime = types.StringValue(audit.LastModifiedTime.Format(time.RFC3339))
	}

	return types.ObjectValue(AuditAttrTypes, map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	})
}
