package view

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

var _ datasource.DataSource = &ViewsDataSource{}
var _ datasource.DataSourceWithConfigure = &ViewsDataSource{}

var DSLListAuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

type ViewsDataSource struct {
	client *client.Client
}

type ViewsDataSourceModel struct {
	Metalake types.String `tfsdk:"metalake"`
	Catalog  types.String `tfsdk:"catalog"`
	Schema   types.String `tfsdk:"schema"`
	Views    types.List   `tfsdk:"views"`
}

func NewViewsDataSource() datasource.DataSource {
	return &ViewsDataSource{}
}

func (d *ViewsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_views"
}

func (d *ViewsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all views within a Gravitino metalake, catalog, and schema.",
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
			"views": schema.ListNestedAttribute{
				Description: "List of views with their details.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The view name.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "The view comment.",
							Computed:    true,
						},
						"view_def": schema.StringAttribute{
							Description: "The SQL view definition.",
							Computed:    true,
						},
						"properties": schema.MapAttribute{
							Description: "Key-value properties for the view.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"audit": schema.ObjectAttribute{
							Description:    "Audit information for the view.",
							Computed:       true,
							AttributeTypes: DSLListAuditAttrTypes,
						},
					},
				},
			},
		},
	}
}

func (ds *ViewsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Invalid provider data", "Expected *client.Client, got unexpected type.")
		return
	}
	ds.client = c
}

func (ds *ViewsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ViewsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	views, err := ds.client.ListViewsDetails(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list views", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(views))
	for i := range views {
		item, d := dslViewListItemToObject(ctx, &views[i])
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		items = []attr.Value{}
	}

	listVal, d := types.ListValue(
		types.ObjectType{AttrTypes: dslViewListItemAttrTypes()},
		items,
	)
	resp.Diagnostics.Append(d...)
	config.Views = listVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func dslViewListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"comment":    types.StringType,
		"view_def":   types.StringType,
		"properties": types.MapType{ElemType: types.StringType},
		"audit":      types.ObjectType{AttrTypes: DSLListAuditAttrTypes},
	}
}

func dslViewListItemToObject(ctx context.Context, v *models.View) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var props types.Map
	if len(v.Properties) > 0 {
		p, d := types.MapValueFrom(ctx, types.StringType, v.Properties)
		diags.Append(d...)
		props = p
	} else {
		props = types.MapNull(types.StringType)
	}

	auditObj, d := dslListAuditToObject(v.Audit)
	diags.Append(d...)

	obj, d := types.ObjectValue(dslViewListItemAttrTypes(), map[string]attr.Value{
		"name":       types.StringValue(v.Name),
		"comment":    types.StringValue(v.Comment),
		"view_def":   types.StringValue(v.ViewDef),
		"properties": props,
		"audit":      auditObj,
	})
	diags.Append(d...)

	return obj, diags
}

func dslListAuditToObject(audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(DSLListAuditAttrTypes), nil
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

	return types.ObjectValue(DSLListAuditAttrTypes, map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	})
}
