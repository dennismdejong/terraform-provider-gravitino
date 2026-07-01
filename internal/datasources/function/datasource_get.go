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

var _ datasource.DataSource = &FunctionDataSource{}
var _ datasource.DataSourceWithConfigure = &FunctionDataSource{}

var dsAuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

type FunctionDataSource struct {
	client *client.Client
}

type FunctionDataSourceModel struct {
	Metalake     types.String `tfsdk:"metalake"`
	Catalog      types.String `tfsdk:"catalog"`
	Schema       types.String `tfsdk:"schema"`
	Name         types.String `tfsdk:"name"`
	Comment      types.String `tfsdk:"comment"`
	FunctionBody types.String `tfsdk:"function_body"`
	Properties   types.Map    `tfsdk:"properties"`
	Audit        types.Object `tfsdk:"audit"`
}

func NewFunctionDataSource() datasource.DataSource {
	return &FunctionDataSource{}
}

func (d *FunctionDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_function"
}

func (d *FunctionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single Gravitino function by name.",
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
			"name": schema.StringAttribute{
				Description: "The function name.",
				Required:    true,
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
				AttributeTypes: dsAuditAttrTypes,
			},
		},
	}
}

func (ds *FunctionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (ds *FunctionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FunctionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	functionResp, err := ds.client.GetFunction(config.Metalake.ValueString(), config.Catalog.ValueString(), config.Schema.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read function", err.Error())
		return
	}

	config.Comment = types.StringValue(functionResp.Function.Comment)
	config.FunctionBody = types.StringValue(functionResp.Function.FunctionBody)

	if len(functionResp.Function.Properties) > 0 {
		props, d := types.MapValueFrom(ctx, types.StringType, functionResp.Function.Properties)
		resp.Diagnostics.Append(d...)
		config.Properties = props
	} else {
		config.Properties = types.MapNull(types.StringType)
	}

	auditObj, d := dsAuditToObject(functionResp.Function.Audit)
	resp.Diagnostics.Append(d...)
	config.Audit = auditObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func dsAuditToObject(audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(dsAuditAttrTypes), nil
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

	return types.ObjectValue(dsAuditAttrTypes, map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	})
}
