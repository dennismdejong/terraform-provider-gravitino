package function

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

var _ datasource.DataSource = &FunctionsDataSource{}
var _ datasource.DataSourceWithConfigure = &FunctionsDataSource{}

var dslAuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

type FunctionsDataSource struct {
	client *client.Client
}

type FunctionsDataSourceModel struct {
	Metalake  types.String `tfsdk:"metalake"`
	Catalog   types.String `tfsdk:"catalog"`
	Schema    types.String `tfsdk:"schema"`
	Functions types.List   `tfsdk:"functions"`
}

func NewFunctionsDataSource() datasource.DataSource {
	return &FunctionsDataSource{}
}

func (d *FunctionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_functions"
}

func (d *FunctionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all functions within a Gravitino metalake, catalog, and schema.",
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
			"functions": schema.ListNestedAttribute{
				Description: "List of functions with their details.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The function name.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "The function comment.",
							Computed:    true,
						},
						"function_body": schema.StringAttribute{
							Description: "The function body.",
							Computed:    true,
						},
						"properties": schema.MapAttribute{
							Description: "Key-value properties for the function.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"audit": schema.ObjectAttribute{
							Description:    "Audit information for the function.",
							Computed:       true,
							AttributeTypes: dslAuditAttrTypes,
						},
					},
				},
			},
		},
	}
}

func (ds *FunctionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (ds *FunctionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FunctionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	functions, err := ds.client.ListFunctionsDetails(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list functions", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(functions))
	for i := range functions {
		item, d := dslFunctionListItemToObject(ctx, &functions[i])
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
		types.ObjectType{AttrTypes: dslFunctionListItemAttrTypes()},
		items,
	)
	resp.Diagnostics.Append(d...)
	config.Functions = listVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func dslFunctionListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"comment":       types.StringType,
		"function_body": types.StringType,
		"properties":    types.MapType{ElemType: types.StringType},
		"audit":         types.ObjectType{AttrTypes: dslAuditAttrTypes},
	}
}

func dslFunctionListItemToObject(ctx context.Context, f *models.Function) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var props types.Map
	if len(f.Properties) > 0 {
		p, d := types.MapValueFrom(ctx, types.StringType, f.Properties)
		diags.Append(d...)
		props = p
	} else {
		props = types.MapNull(types.StringType)
	}

	auditObj, d := dslAuditToObject(f.Audit)
	diags.Append(d...)

	obj, d := types.ObjectValue(dslFunctionListItemAttrTypes(), map[string]attr.Value{
		"name":          types.StringValue(f.Name),
		"comment":       types.StringValue(f.Comment),
		"function_body": types.StringValue(f.FunctionBody),
		"properties":    props,
		"audit":         auditObj,
	})
	diags.Append(d...)

	return obj, diags
}

func dslAuditToObject(audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(dslAuditAttrTypes), nil
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

	return types.ObjectValue(dslAuditAttrTypes, map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	})
}
