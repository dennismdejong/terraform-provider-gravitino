package job

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
)

var _ resource.Resource = &JobResource{}
var _ resource.ResourceWithImportState = &JobResource{}
var _ resource.ResourceWithConfigure = &JobResource{}

type JobResource struct {
	client *client.Client
}

func New() resource.Resource {
	return &JobResource{}
}

func (r *JobResource) SetClient(c *client.Client) {
	r.client = c
}

type JobResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Metalake   types.String `tfsdk:"metalake"`
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Parameters types.Map    `tfsdk:"parameters"`
	Schedule   types.String `tfsdk:"schedule"`
	Status     types.String `tfsdk:"status"`
	Audit      types.Object `tfsdk:"audit"`
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (r *JobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JobResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "gravitino_job"
}

func (r *JobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier in the format 'metalake.job_name'.",
			},
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The job name.",
			},
			"template": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The job template name.",
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value parameters for the job.",
			},
			"schedule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The job schedule (cron expression).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the job.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the job.",
			},
		},
	}
}

func (r *JobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.JobCreateRequest{
		Name:       plan.Name.ValueString(),
		Template:   plan.Template.ValueString(),
		Parameters: mapFromTFInterface(plan.Parameters),
		Schedule:   plan.Schedule.ValueString(),
	}

	result, err := r.client.CreateJob(plan.Metalake.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create job", err.Error())
		return
	}

	setStateFromJob(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), &result.Job, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *JobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetJob(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read job", err.Error())
		return
	}

	setStateFromJob(ctx, &resp.Diagnostics, state.Metalake.ValueString(), &result.Job, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state JobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setStateFromJob(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), nil, &plan)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *JobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state JobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteJob(state.Metalake.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete job", err.Error())
		return
	}
}

func (r *JobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idx := strings.LastIndex(req.ID, ".")
	if idx == -1 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected 'metalake.job_name', got: %s", req.ID),
		)
		return
	}

	metalake := req.ID[:idx]
	name := req.ID[idx+1:]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metalake"), metalake)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setStateFromJob(ctx context.Context, diags *diag.Diagnostics, metalake string, job *models.Job, model *JobResourceModel) {
	if job != nil {
		model.Name = types.StringValue(job.Name)
		model.Template = types.StringValue(job.Template)
		model.Schedule = types.StringValue(job.Schedule)
		model.Status = types.StringValue(job.Status)
		model.ID = types.StringValue(metalake + "." + job.Name)

		props, d := types.MapValueFrom(ctx, types.StringType, strMapFromInterface(job.Parameters))
		diags.Append(d...)
		if !diags.HasError() {
			model.Parameters = props
		}

		if job.Audit != nil {
			auditObj, d := auditToObjectValue(ctx, job.Audit)
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

func mapFromTFInterface(m types.Map) map[string]interface{} {
	result := make(map[string]interface{})
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

func strMapFromInterface(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			result[k] = strVal
		}
	}
	return result
}
