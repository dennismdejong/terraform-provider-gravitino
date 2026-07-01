package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PolicyResource{}
var _ resource.ResourceWithImportState = &PolicyResource{}
var _ resource.ResourceWithConfigure = &PolicyResource{}

type PolicyResource struct {
	client *client.Client
}

func New() resource.Resource {
	return &PolicyResource{}
}

func (r *PolicyResource) SetClient(c *client.Client) {
	r.client = c
}

type PolicyResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Metalake     types.String `tfsdk:"metalake"`
	ResourceType types.String `tfsdk:"resource_type"`
	Resource     types.String `tfsdk:"resource"`
	Name         types.String `tfsdk:"name"`
	Condition    types.String `tfsdk:"condition"`
	Effect       types.String `tfsdk:"effect"`
	Actions      types.List   `tfsdk:"actions"`
	Subjects     types.List   `tfsdk:"subjects"`
	Object       types.String `tfsdk:"object"`
	Properties   types.Map    `tfsdk:"properties"`
	Audit        types.Object `tfsdk:"audit"`
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (r *PolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *PolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "gravitino_policy"
}

func (r *PolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier in the format 'metalake.resource_type.resource.policy_name'.",
			},
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The metadata object resource type. Must be one of: metalakes, catalogs, schemas, tables, filesets, topics.",
				Validators: []validator.String{
					stringvalidator.OneOf("metalakes", "catalogs", "schemas", "tables", "filesets", "topics"),
				},
			},
			"resource": schema.StringAttribute{
				Required:    true,
				Description: "The metadata object name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The policy name.",
			},
			"condition": schema.StringAttribute{
				Optional:    true,
				Description: "The policy condition expression.",
			},
			"effect": schema.StringAttribute{
				Required:    true,
				Description: "The policy effect. Must be one of: allow, deny.",
				Validators: []validator.String{
					stringvalidator.OneOf("allow", "deny"),
				},
			},
			"actions": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The list of actions this policy applies to.",
			},
			"subjects": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The list of subjects this policy applies to.",
			},
			"object": schema.StringAttribute{
				Optional:    true,
				Description: "The policy target object.",
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value properties for the policy.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the policy.",
			},
		},
	}
}

func (r *PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating policy", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})

	createReq := &models.PolicyCreateRequest{
		Name:       plan.Name.ValueString(),
		Effect:     plan.Effect.ValueString(),
		Condition:  plan.Condition.ValueString(),
		Object:     plan.Object.ValueString(),
		Actions:    listFromTF(plan.Actions),
		Subjects:   listFromTF(plan.Subjects),
		Properties: mapFromTF(plan.Properties),
	}

	metalake := plan.Metalake.ValueString()
	resourceType := plan.ResourceType.ValueString()
	resourceName := plan.Resource.ValueString()

	result, err := r.client.CreatePolicy(metalake, resourceType, resourceName, createReq)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("creating policy", plan.Name.ValueString(), err)...)
		return
	}

	setStateFromPolicy(ctx, &resp.Diagnostics, metalake, resourceType, resourceName, &result.Policy, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Created policy", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading policy", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	metalake := state.Metalake.ValueString()
	resourceType := state.ResourceType.ValueString()
	resourceName := state.Resource.ValueString()

	result, err := r.client.GetPolicy(metalake, resourceType, resourceName, state.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(client.NewResourceError("reading policy", state.Name.ValueString(), err)...)
		return
	}

	setStateFromPolicy(ctx, &resp.Diagnostics, metalake, resourceType, resourceName, &result.Policy, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating policy", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	var updates []interface{}

	metalake := plan.Metalake.ValueString()
	resourceType := plan.ResourceType.ValueString()
	resourceName := plan.Resource.ValueString()

	if !plan.Name.Equal(state.Name) {
		updates = append(updates, models.NewRenamePolicyRequest(plan.Name.ValueString()))
		state.Name = plan.Name
	}

	oldProps := mapFromTF(state.Properties)
	newProps := mapFromTF(plan.Properties)

	for k, v := range newProps {
		oldV, exists := oldProps[k]
		if !exists || oldV != v {
			updates = append(updates, models.NewSetPolicyPropertyRequest(k, v))
		}
	}

	for k := range oldProps {
		if _, exists := newProps[k]; !exists {
			updates = append(updates, models.NewRemovePolicyPropertyRequest(k))
		}
	}

	if len(updates) > 0 {
		result, err := r.client.UpdatePolicy(metalake, resourceType, resourceName, state.Name.ValueString(), updates)
		if err != nil {
			resp.Diagnostics.Append(client.NewResourceError("updating policy", state.Name.ValueString(), err)...)
			return
		}
		setStateFromPolicy(ctx, &resp.Diagnostics, metalake, resourceType, resourceName, &result.Policy, &plan)
	} else {
		setStateFromPolicy(ctx, &resp.Diagnostics, metalake, resourceType, resourceName, nil, &plan)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Updated policy", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting policy", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	_, err := r.client.DeletePolicy(
		state.Metalake.ValueString(),
		state.ResourceType.ValueString(),
		state.Resource.ValueString(),
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("deleting policy", state.Name.ValueString(), err)...)
		return
	}

	tflog.Debug(ctx, "Deleted policy", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
}

func (r *PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ".", 4)
	if len(parts) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected 'metalake.resource_type.resource.policy_name', got: %s", req.ID),
		)
		return
	}

	metalake, resourceType, resourceName, name := parts[0], parts[1], parts[2], parts[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metalake"), metalake)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_type"), resourceType)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource"), resourceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setStateFromPolicy(ctx context.Context, diags *diag.Diagnostics, metalake, resourceType, resourceName string, policy *models.Policy, model *PolicyResourceModel) {
	if policy != nil {
		model.Name = types.StringValue(policy.Name)
		model.Effect = types.StringValue(policy.Effect)
		model.Condition = types.StringValue(policy.Condition)
		model.Object = types.StringValue(policy.Object)
		model.ID = types.StringValue(metalake + "." + resourceType + "." + resourceName + "." + policy.Name)

		actions, d := types.ListValueFrom(ctx, types.StringType, policy.Actions)
		diags.Append(d...)
		if !diags.HasError() {
			model.Actions = actions
		}

		subjects, d := types.ListValueFrom(ctx, types.StringType, policy.Subjects)
		diags.Append(d...)
		if !diags.HasError() {
			model.Subjects = subjects
		}

		props, d := types.MapValueFrom(ctx, types.StringType, policy.Properties)
		diags.Append(d...)
		if !diags.HasError() {
			model.Properties = props
		}

		if policy.Audit != nil {
			auditObj, d := auditToObjectValue(ctx, policy.Audit)
			diags.Append(d...)
			if !diags.HasError() {
				model.Audit = auditObj
			}
		}
	} else {
		model.ID = types.StringValue(metalake + "." + resourceType + "." + resourceName + "." + model.Name.ValueString())
	}

	model.Metalake = types.StringValue(metalake)
	model.ResourceType = types.StringValue(resourceType)
	model.Resource = types.StringValue(resourceName)
}

func auditToObjectValue(ctx context.Context, audit *models.Audit) (types.Object, diag.Diagnostics) {
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

func mapFromTF(m types.Map) map[string]string {
	result := make(map[string]string)
	if m.IsNull() || m.IsUnknown() {
		return result
	}
	for k, v := range m.Elements() {
		if strVal, ok := v.(types.String); ok {
			result[k] = strVal.ValueString()
		}
	}
	return result
}

func listFromTF(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	var result []string
	for _, v := range l.Elements() {
		if strVal, ok := v.(types.String); ok {
			result = append(result, strVal.ValueString())
		}
	}
	return result
}
