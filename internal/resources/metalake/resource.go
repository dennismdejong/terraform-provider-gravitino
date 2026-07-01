package metalake

import (
	"context"
	"fmt"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*MetalakeResource)(nil)
var _ resource.ResourceWithImportState = (*MetalakeResource)(nil)

type MetalakeResource struct {
	client *client.Client
}

type MetalakeResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      *AuditModel  `tfsdk:"audit"`
}

type AuditModel struct {
	Creator          types.String `tfsdk:"creator"`
	CreateTime       types.String `tfsdk:"create_time"`
	LastModifier     types.String `tfsdk:"last_modifier"`
	LastModifiedTime types.String `tfsdk:"last_modified_time"`
}

func NewMetalakeResource() resource.Resource {
	return &MetalakeResource{}
}

func (r *MetalakeResource) SetClient(c *client.Client) {
	r.client = c
}

func (r *MetalakeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metalake"
}

func (r *MetalakeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"audit": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"creator": schema.StringAttribute{
						Computed: true,
					},
					"create_time": schema.StringAttribute{
						Computed: true,
					},
					"last_modifier": schema.StringAttribute{
						Computed: true,
					},
					"last_modified_time": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *MetalakeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = cli
}

func (r *MetalakeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MetalakeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := mapToProperties(ctx, plan.Properties, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.MetalakeCreateRequest{
		Name:       plan.Name.ValueString(),
		Comment:    plan.Comment.ValueString(),
		Properties: props,
	}

	result, err := r.client.CreateMetalake(createReq)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("creating metalake", plan.Name.ValueString(), err)...)
		return
	}

	plan.ID = plan.Name

	metalakeToState(&result.Metalake, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *MetalakeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MetalakeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetMetalake(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("reading metalake", state.Name.ValueString(), err)...)
		return
	}

	state.ID = state.Name

	metalakeToState(&result.Metalake, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MetalakeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MetalakeResourceModel
	var state MetalakeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updates []interface{}

	if !plan.Name.Equal(state.Name) {
		updates = append(updates, models.NewRenameMetalakeRequest(plan.Name.ValueString()))
	}

	if !plan.Comment.Equal(state.Comment) {
		newComment := plan.Comment.ValueString()
		updates = append(updates, models.NewUpdateMetalakeCommentRequest(newComment))
	}

	oldProps := mapToProperties(ctx, state.Properties, &resp.Diagnostics)
	newProps := mapToProperties(ctx, plan.Properties, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	for key, oldVal := range oldProps {
		newVal, exists := newProps[key]
		if !exists {
			updates = append(updates, models.NewRemoveMetalakePropertyRequest(key))
		} else if oldVal != newVal {
			updates = append(updates, models.NewSetMetalakePropertyRequest(key, newVal))
		}
	}

	for key, newVal := range newProps {
		if _, exists := oldProps[key]; !exists {
			updates = append(updates, models.NewSetMetalakePropertyRequest(key, newVal))
		}
	}

	result, err := r.client.UpdateMetalake(state.Name.ValueString(), updates)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("updating metalake", state.Name.ValueString(), err)...)
		return
	}

	plan.ID = plan.Name

	metalakeToState(&result.Metalake, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *MetalakeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MetalakeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DropMetalake(state.Name.ValueString(), true)
	if err != nil {
		resp.Diagnostics.Append(client.NewResourceError("deleting metalake", state.Name.ValueString(), err)...)
		return
	}
}

func (r *MetalakeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func mapToProperties(ctx context.Context, m types.Map, diags *diag.Diagnostics) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	result := make(map[string]string, len(m.Elements()))
	d := m.ElementsAs(ctx, &result, false)
	if d.HasError() {
		*diags = append(*diags, d...)
		return nil
	}
	return result
}

func propertiesToMap(ctx context.Context, props map[string]string, diags *diag.Diagnostics) types.Map {
	if len(props) == 0 {
		return types.MapNull(types.StringType)
	}
	result, d := types.MapValueFrom(ctx, types.StringType, props)
	if d.HasError() {
		*diags = append(*diags, d...)
		return types.MapNull(types.StringType)
	}
	return result
}

func timeToString(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func metalakeToState(m *models.Metalake, state *MetalakeResourceModel, diags *diag.Diagnostics) {
	state.Name = types.StringValue(m.Name)
	state.Comment = types.StringValue(m.Comment)
	state.Properties = propertiesToMap(context.Background(), m.Properties, diags)

	if m.Audit != nil {
		state.Audit = &AuditModel{
			Creator:          types.StringValue(m.Audit.Creator),
			CreateTime:       timeToString(m.Audit.CreateTime),
			LastModifier:     types.StringValue(m.Audit.LastModifier),
			LastModifiedTime: timeToString(m.Audit.LastModifiedTime),
		}
	} else {
		state.Audit = nil
	}
}
