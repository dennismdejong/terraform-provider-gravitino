package job_template

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

var _ resource.Resource = &JobTemplateResource{}
var _ resource.ResourceWithImportState = &JobTemplateResource{}
var _ resource.ResourceWithConfigure = &JobTemplateResource{}

type JobTemplateResource struct {
	client *client.Client
}

func NewJobTemplateResource() resource.Resource {
	return &JobTemplateResource{}
}

func (r *JobTemplateResource) SetClient(c *client.Client) {
	r.client = c
}

type JobTemplateResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Metalake   types.String `tfsdk:"metalake"`
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Parameters types.Map    `tfsdk:"parameters"`
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

func (r *JobTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JobTemplateResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "gravitino_job_template"
}

func (r *JobTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier in the format 'metalake.template_name'.",
			},
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The job template name.",
			},
			"template": schema.StringAttribute{
				Required:    true,
				Description: "The template definition.",
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value parameters for the job template.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A comment or description for the job template.",
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value properties for the job template.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the job template.",
			},
		},
	}
}

func (r *JobTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JobTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating job template", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})

	createReq := &models.JobTemplateCreateRequest{
		Name:       plan.Name.ValueString(),
		Template:   plan.Template.ValueString(),
		Parameters: mapFromTF(plan.Parameters),
		Comment:    plan.Comment.ValueString(),
		Properties: mapFromTF(plan.Properties),
	}

	result, err := r.client.RegisterJobTemplate(plan.Metalake.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("creating job template", plan.Name.ValueString(), err)...)
		return
	}

	setStateFromJobTemplate(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), &result.JobTemplate, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Created job template", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *JobTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JobTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading job template", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	result, err := r.client.GetJobTemplate(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(client.NewResourceError("reading job template", state.Name.ValueString(), err)...)
		return
	}

	setStateFromJobTemplate(ctx, &resp.Diagnostics, state.Metalake.ValueString(), &result.JobTemplate, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *JobTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state JobTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating job template", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	needsUpdate := !plan.Template.Equal(state.Template) ||
		!plan.Comment.Equal(state.Comment) ||
		!mapsEqual(plan.Parameters, state.Parameters) ||
		!mapsEqual(plan.Properties, state.Properties)

	if needsUpdate {
		updateReq := &models.JobTemplateCreateRequest{
			Name:       plan.Name.ValueString(),
			Template:   plan.Template.ValueString(),
			Parameters: mapFromTF(plan.Parameters),
			Comment:    plan.Comment.ValueString(),
			Properties: mapFromTF(plan.Properties),
		}

		result, err := r.client.UpdateJobTemplate(plan.Metalake.ValueString(), plan.Name.ValueString(), updateReq)
		if err != nil {
			resp.Diagnostics.Append(client.NewResourceError("updating job template", plan.Name.ValueString(), err)...)
			return
		}
		setStateFromJobTemplate(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), &result.JobTemplate, &plan)
	} else {
		setStateFromJobTemplate(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), nil, &plan)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "Updated job template", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
}

func (r *JobTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JobTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting job template", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

	_, err := r.client.DeleteJobTemplate(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("deleting job template", state.Name.ValueString(), err)...)
		return
	}

	tflog.Debug(ctx, "Deleted job template", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
}

func (r *JobTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idx := strings.LastIndex(req.ID, ".")
	if idx == -1 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected 'metalake.template_name', got: %s", req.ID),
		)
		return
	}

	metalake := req.ID[:idx]
	name := req.ID[idx+1:]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metalake"), metalake)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setStateFromJobTemplate(ctx context.Context, diags *diag.Diagnostics, metalake string, jt *models.JobTemplate, model *JobTemplateResourceModel) {
	if jt != nil {
		model.Name = types.StringValue(jt.Name)
		model.Template = types.StringValue(jt.Template)
		model.Comment = types.StringValue(jt.Comment)
		model.ID = types.StringValue(metalake + "." + jt.Name)

		params, d := types.MapValueFrom(ctx, types.StringType, jt.Parameters)
		diags.Append(d...)
		if !diags.HasError() {
			model.Parameters = params
		}

		props, d := types.MapValueFrom(ctx, types.StringType, jt.Properties)
		diags.Append(d...)
		if !diags.HasError() {
			model.Properties = props
		}

		if jt.Audit != nil {
			auditObj, d := auditToObjectValue(ctx, jt.Audit)
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

	attrs := map[string]attr.Value{
		"creator":       types.StringValue(audit.Creator),
		"last_modifier": types.StringValue(audit.LastModifier),
	}

	if audit.CreateTime != nil {
		attrs["create_time"] = types.StringValue(audit.CreateTime.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		attrs["create_time"] = types.StringNull()
	}
	if audit.LastModifiedTime != nil {
		attrs["last_modified_time"] = types.StringValue(audit.LastModifiedTime.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		attrs["last_modified_time"] = types.StringNull()
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

func mapsEqual(a, b types.Map) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	if len(a.Elements()) != len(b.Elements()) {
		return false
	}
	for k, va := range a.Elements() {
		vb, ok := b.Elements()[k]
		if !ok {
			return false
		}
		sa, okA := va.(types.String)
		sb, okB := vb.(types.String)
		if !okA || !okB {
			return false
		}
		if !sa.Equal(sb) {
			return false
		}
	}
	return true
}
