package model

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

var _ datasource.DataSource = &ModelsDataSource{}
var _ datasource.DataSourceWithConfigure = &ModelsDataSource{}

var dslAuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

type ModelsDataSource struct {
	client *client.Client
}

type ModelsDataSourceModel struct {
	Metalake types.String `tfsdk:"metalake"`
	Catalog  types.String `tfsdk:"catalog"`
	Schema   types.String `tfsdk:"schema"`
	Models   types.List   `tfsdk:"models"`
}

func NewModelsDataSource() datasource.DataSource {
	return &ModelsDataSource{}
}

func (d *ModelsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_models"
}

func (d *ModelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all models within a Gravitino metalake, catalog, and schema.",
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
			"models": schema.ListNestedAttribute{
				Description: "List of models with their details.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The model name.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "The model comment.",
							Computed:    true,
						},
						"model_uri": schema.StringAttribute{
							Description: "The URI of the model artifact.",
							Computed:    true,
						},
						"properties": schema.MapAttribute{
							Description: "Key-value properties for the model.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"audit": schema.ObjectAttribute{
							Description:    "Audit information for the model.",
							Computed:       true,
							AttributeTypes: dslAuditAttrTypes,
						},
					},
				},
			},
		},
	}
}

func (d *ModelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ModelsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mods, err := d.client.ListModelsDetails(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list models", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(mods))
	for i := range mods {
		item, itemDiags := modelListItemToObject(ctx, &mods[i])
		resp.Diagnostics.Append(itemDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		items = []attr.Value{}
	}

	listVal, listDiags := types.ListValue(
		types.ObjectType{AttrTypes: modelListItemAttrTypes()},
		items,
	)
	resp.Diagnostics.Append(listDiags...)
	config.Models = listVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func modelListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"comment":    types.StringType,
		"model_uri":  types.StringType,
		"properties": types.MapType{ElemType: types.StringType},
		"audit":      types.ObjectType{AttrTypes: dslAuditAttrTypes},
	}
}

func modelListItemToObject(ctx context.Context, m *models.Model) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var props types.Map
	if len(m.Properties) > 0 {
		p, pDiags := types.MapValueFrom(ctx, types.StringType, m.Properties)
		diags.Append(pDiags...)
		props = p
	} else {
		props = types.MapNull(types.StringType)
	}

	ao, aDiags := listAuditToObject(m.Audit)
	diags.Append(aDiags...)

	obj, oDiags := types.ObjectValue(modelListItemAttrTypes(), map[string]attr.Value{
		"name":       types.StringValue(m.Name),
		"comment":    types.StringValue(m.Comment),
		"model_uri":  types.StringValue(m.ModelURI),
		"properties": props,
		"audit":      ao,
	})
	diags.Append(oDiags...)

	return obj, diags
}

func listAuditToObject(a *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if a == nil {
		return types.ObjectNull(dslAuditAttrTypes), nil
	}

	creator := types.StringNull()
	if a.Creator != "" {
		creator = types.StringValue(a.Creator)
	}

	createTime := types.StringNull()
	if a.CreateTime != nil {
		createTime = types.StringValue(a.CreateTime.Format(time.RFC3339))
	}

	lastModifier := types.StringNull()
	if a.LastModifier != "" {
		lastModifier = types.StringValue(a.LastModifier)
	}

	lastModifiedTime := types.StringNull()
	if a.LastModifiedTime != nil {
		lastModifiedTime = types.StringValue(a.LastModifiedTime.Format(time.RFC3339))
	}

	return types.ObjectValue(dslAuditAttrTypes, map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	})
}
