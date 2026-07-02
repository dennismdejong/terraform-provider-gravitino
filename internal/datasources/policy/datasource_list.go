package policy

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ datasource.DataSource = &PoliciesDataSource{}
var _ datasource.DataSourceWithConfigure = &PoliciesDataSource{}

type PoliciesDataSource struct {
	client *client.Client
}

func NewListDataSource() datasource.DataSource {
	return &PoliciesDataSource{}
}

func (d *PoliciesDataSource) SetClient(c *client.Client) {
	d.client = c
}

type PoliciesDataSourceModel struct {
	Metalake     types.String `tfsdk:"metalake"`
	ResourceType types.String `tfsdk:"resource_type"`
	Resource     types.String `tfsdk:"resource"`
	Policies     types.List   `tfsdk:"policies"`
}

type policyItemModel struct {
	Name       types.String `tfsdk:"name"`
	Condition  types.String `tfsdk:"condition"`
	Effect     types.String `tfsdk:"effect"`
	Actions    types.List   `tfsdk:"actions"`
	Subjects   types.List   `tfsdk:"subjects"`
	Object     types.String `tfsdk:"object"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      types.Object `tfsdk:"audit"`
}

var PolicyItemAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"condition":  types.StringType,
	"effect":     types.StringType,
	"actions":    types.ListType{ElemType: types.StringType},
	"subjects":   types.ListType{ElemType: types.StringType},
	"object":     types.StringType,
	"properties": types.MapType{ElemType: types.StringType},
	"audit":      types.ObjectType{AttrTypes: AuditAttrTypes},
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (d *PoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoliciesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_policies"
}

func (d *PoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The metadata object resource type. Must be one of: metalakes, catalogs, schemas, tables, filesets, topics.",
				Validators: []validator.String{
					stringvalidator.OneOf("metalakes", "catalogs", "schemas", "tables", "filesets", "topics", "models"),
				},
			},
			"resource": schema.StringAttribute{
				Required:    true,
				Description: "The metadata object name.",
			},
			"policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The policy name.",
						},
						"condition": schema.StringAttribute{
							Computed:    true,
							Description: "The policy condition expression.",
						},
						"effect": schema.StringAttribute{
							Computed:    true,
							Description: "The policy effect.",
						},
						"actions": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The list of actions.",
						},
						"subjects": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The list of subjects.",
						},
						"object": schema.StringAttribute{
							Computed:    true,
							Description: "The policy target object.",
						},
						"properties": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The policy properties.",
						},
						"audit": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: AuditAttrTypes,
							Description:    "Audit information for the policy.",
						},
					},
				},
			},
		},
	}
}

func (d *PoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListPolicies(config.Metalake.ValueString(), config.ResourceType.ValueString(), config.Resource.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list policies", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(result.Policies))
	for _, p := range result.Policies {
		policy := p
		item := policyToItemModel(ctx, &policy, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if item == nil {
			continue
		}
		obj, objDiags := types.ObjectValueFrom(ctx, PolicyItemAttrTypes, item)
		resp.Diagnostics.Append(objDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, obj)
	}

	policiesList, listDiags := types.ListValue(types.ObjectType{AttrTypes: PolicyItemAttrTypes}, items)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Policies = policiesList
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func policyToItemModel(ctx context.Context, p *models.Policy, diags *diag.Diagnostics) *policyItemModel {
	if p == nil {
		return nil
	}

	item := &policyItemModel{
		Name:      types.StringValue(p.Name),
		Condition: types.StringValue(p.Condition),
		Effect:    types.StringValue(p.Effect),
		Object:    types.StringValue(p.Object),
	}

	actions, d := types.ListValueFrom(ctx, types.StringType, p.Actions)
	if d.HasError() {
		return nil
	}
	item.Actions = actions

	subjects, d := types.ListValueFrom(ctx, types.StringType, p.Subjects)
	if d.HasError() {
		return nil
	}
	item.Subjects = subjects

	props, d := types.MapValueFrom(ctx, types.StringType, p.Properties)
	if d.HasError() {
		return nil
	}
	item.Properties = props

	auditObj, d := auditToObjectValueForDS(ctx, p.Audit)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}
	item.Audit = auditObj

	return item
}

func auditToObjectValueForDS(ctx context.Context, audit *models.Audit) (basetypes.ObjectValue, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(AuditAttrTypes), nil
	}

	creator := types.StringValue(audit.Creator)
	lastModifier := types.StringValue(audit.LastModifier)

	var createTime, lastModifiedTime types.String
	if audit.CreateTime != nil {
		createTime = types.StringValue(audit.CreateTime.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		createTime = types.StringNull()
	}
	if audit.LastModifiedTime != nil {
		lastModifiedTime = types.StringValue(audit.LastModifiedTime.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		lastModifiedTime = types.StringNull()
	}

	attrs := map[string]attr.Value{
		"creator":            creator,
		"create_time":        createTime,
		"last_modifier":      lastModifier,
		"last_modified_time": lastModifiedTime,
	}

	return types.ObjectValue(AuditAttrTypes, attrs)
}
