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

var _ datasource.DataSource = &PartitionsDataSource{}
var _ datasource.DataSourceWithConfigure = &PartitionsDataSource{}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

type PartitionsDataSource struct {
	client *client.Client
}

type PartitionsDataSourceModel struct {
	Metalake   types.String `tfsdk:"metalake"`
	Catalog    types.String `tfsdk:"catalog"`
	Schema     types.String `tfsdk:"schema"`
	Table      types.String `tfsdk:"table"`
	Partitions types.List   `tfsdk:"partitions"`
}

func NewPartitionsDataSource() datasource.DataSource {
	return &PartitionsDataSource{}
}

func (d *PartitionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_partitions"
}

func (d *PartitionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all partitions within a Gravitino table.",
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
			"partitions": schema.ListNestedAttribute{
				Description: "List of partitions with their details.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The partition name.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *PartitionsDataSource) SetClient(c *client.Client) {
	d.client = c
}

func (d *PartitionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PartitionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PartitionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	partitions, err := d.client.ListPartitionsDetails(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString(), config.Table.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list partitions", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(partitions))
	for i := range partitions {
		item, d := dslPartitionListItemToObject(ctx, &partitions[i])
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		items = []attr.Value{}
	}

	listVal, listDiags := types.ListValue(
		types.ObjectType{AttrTypes: PartitionItemAttrTypes},
		items,
	)
	resp.Diagnostics.Append(listDiags...)
	config.Partitions = listVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

var PartitionItemAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"properties": types.MapType{ElemType: types.StringType},
	"audit":      types.ObjectType{AttrTypes: AuditAttrTypes},
}

func dslPartitionListItemToObject(ctx context.Context, p *models.Partition) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var props types.Map
	if len(p.Properties) > 0 {
		mp, d := types.MapValueFrom(ctx, types.StringType, p.Properties)
		diags.Append(d...)
		props = mp
	} else {
		props = types.MapNull(types.StringType)
	}

	auditObj, d := dslAuditToObject(p.Audit)
	diags.Append(d...)

	obj, d := types.ObjectValue(PartitionItemAttrTypes, map[string]attr.Value{
		"name":       types.StringValue(p.Name),
		"properties": props,
		"audit":      auditObj,
	})
	diags.Append(d...)

	return obj, diags
}

func dslAuditToObject(audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
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
