package tag

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}
var _ resource.ResourceWithConfigure = &TagResource{}

type TagResource struct {
	client *client.Client
}

func New() resource.Resource {
	return &TagResource{}
}

func (r *TagResource) SetClient(c *client.Client) {
	r.client = c
}

type TagResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Metalake   types.String `tfsdk:"metalake"`
	Name       types.String `tfsdk:"name"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      types.Object `tfsdk:"audit"`
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (r *TagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TagResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "gravitino_tag"
}

func (r *TagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier in the format 'metalake.tag_name'.",
			},
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The tag name.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A comment or description for the tag.",
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value properties for the tag.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the tag.",
			},
		},
	}
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating tag", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})

	createReq := &models.TagCreateRequest{
		Name:       plan.Name.ValueString(),
		Comment:    plan.Comment.ValueString(),
		Properties: mapFromTF(plan.Properties),
	}

	result, err := r.client.CreateTag(plan.Metalake.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("creating tag", plan.Name.ValueString(), err)...)
		return
	}

	setStateFromTag(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), &result.Tag, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Created tag", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading tag", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	result, err := r.client.GetTag(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(client.NewResourceError("reading tag", state.Name.ValueString(), err)...)
		return
	}

	setStateFromTag(ctx, &resp.Diagnostics, state.Metalake.ValueString(), &result.Tag, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating tag", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	var updates []interface{}

	if !plan.Name.Equal(state.Name) {
		updates = append(updates, models.NewRenameTagRequest(plan.Name.ValueString()))
		state.Name = plan.Name
	}

	if !plan.Comment.Equal(state.Comment) {
		updates = append(updates, models.NewUpdateTagCommentRequest(plan.Comment.ValueString()))
		state.Comment = plan.Comment
	}

	oldProps := mapFromTF(state.Properties)
	newProps := mapFromTF(plan.Properties)

	for k, v := range newProps {
		oldV, exists := oldProps[k]
		if !exists || oldV != v {
			updates = append(updates, models.NewSetTagPropertyRequest(k, v))
		}
	}

	for k := range oldProps {
		if _, exists := newProps[k]; !exists {
			updates = append(updates, models.NewRemoveTagPropertyRequest(k))
		}
	}

	if len(updates) > 0 {
		result, err := r.client.UpdateTag(state.Metalake.ValueString(), state.Name.ValueString(), updates)
		if err != nil {
			resp.Diagnostics.Append(client.NewResourceError("updating tag", state.Name.ValueString(), err)...)
			return
		}
		setStateFromTag(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), &result.Tag, &plan)
	} else {
		setStateFromTag(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), nil, &plan)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Updated tag", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting tag", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	_, err := r.client.DeleteTag(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("deleting tag", state.Name.ValueString(), err)...)
		return
	}

	tflog.Debug(ctx, "Deleted tag", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idx := strings.LastIndex(req.ID, ".")
	if idx == -1 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected 'metalake.tag_name', got: %s", req.ID),
		)
		return
	}

	metalake := req.ID[:idx]
	name := req.ID[idx+1:]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metalake"), metalake)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setStateFromTag(ctx context.Context, diags *diag.Diagnostics, metalake string, tag *models.Tag, model *TagResourceModel) {
	if tag != nil {
		model.Name = types.StringValue(tag.Name)
		model.Comment = types.StringValue(tag.Comment)
		model.ID = types.StringValue(metalake + "." + tag.Name)

		props, d := types.MapValueFrom(ctx, types.StringType, tag.Properties)
		diags.Append(d...)
		if !diags.HasError() {
			model.Properties = props
		}

		if tag.Audit != nil {
			auditObj, d := auditToObjectValue(ctx, tag.Audit)
			diags.Append(d...)
			if !diags.HasError() {
				model.Audit = auditObj
			}
		}
	} else {
		model.ID = types.StringValue(metalake + "." + model.Name.ValueString())
	}

	model.Metalake = types.StringValue(metalake)
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
