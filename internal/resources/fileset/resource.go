package fileset

import (
	"context"
	"fmt"
	"strings"
	"time"

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
)

var _ resource.Resource = &FilesetResource{}
var _ resource.ResourceWithImportState = &FilesetResource{}
var _ resource.ResourceWithConfigure = &FilesetResource{}

type FilesetResource struct {
	client *client.Client
}

func New() resource.Resource {
	return &FilesetResource{}
}

func (r *FilesetResource) SetClient(c *client.Client) {
	r.client = c
}

type FilesetResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Metalake        types.String `tfsdk:"metalake"`
	Catalog         types.String `tfsdk:"catalog"`
	Schema          types.String `tfsdk:"schema"`
	Name            types.String `tfsdk:"name"`
	Comment         types.String `tfsdk:"comment"`
	Type            types.String `tfsdk:"type"`
	StorageLocation types.String `tfsdk:"storage_location"`
	Properties      types.Map    `tfsdk:"properties"`
	Audit           types.Object `tfsdk:"audit"`
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (r *FilesetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesetResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "gravitino_fileset"
}

func (r *FilesetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier in the format 'metalake.catalog.schema.fileset'.",
			},
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"catalog": schema.StringAttribute{
				Required:    true,
				Description: "The catalog name.",
			},
			"schema": schema.StringAttribute{
				Required:    true,
				Description: "The schema name.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The fileset name.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A comment or description for the fileset.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The fileset type. Must be one of: managed, external.",
				Validators: []validator.String{
					stringvalidator.OneOf("managed", "external"),
				},
			},
			"storage_location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The storage location for the fileset.",
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of key-value properties for the fileset.",
			},
			"audit": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: AuditAttrTypes,
				Description:    "Audit information for the fileset.",
			},
		},
	}
}

func (r *FilesetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FilesetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.FilesetCreateRequest{
		Name:            plan.Name.ValueString(),
		Comment:         plan.Comment.ValueString(),
		Type:            plan.Type.ValueString(),
		StorageLocation: plan.StorageLocation.ValueString(),
		Properties:      mapFromTF(plan.Properties),
	}

	result, err := r.client.CreateFileset(plan.Metalake.ValueString(), plan.Catalog.ValueString(), plan.Schema.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create fileset", err.Error())
		return
	}

	setStateFromFileset(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), plan.Catalog.ValueString(), plan.Schema.ValueString(), &result.Fileset, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FilesetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FilesetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetFileset(state.Metalake.ValueString(), state.Catalog.ValueString(), state.Schema.ValueString(), state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read fileset", err.Error())
		return
	}

	setStateFromFileset(ctx, &resp.Diagnostics, state.Metalake.ValueString(), state.Catalog.ValueString(), state.Schema.ValueString(), &result.Fileset, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *FilesetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FilesetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updates []interface{}

	if !plan.Name.Equal(state.Name) {
		updates = append(updates, models.NewRenameFilesetRequest(plan.Name.ValueString()))
		state.Name = plan.Name
	}

	if !plan.Comment.Equal(state.Comment) {
		updates = append(updates, models.NewUpdateFilesetCommentRequest(plan.Comment.ValueString()))
		state.Comment = plan.Comment
	}

	oldProps := mapFromTF(state.Properties)
	newProps := mapFromTF(plan.Properties)

	for k, v := range newProps {
		oldV, exists := oldProps[k]
		if !exists || oldV != v {
			updates = append(updates, models.NewSetFilesetPropertyRequest(k, v))
		}
	}

	for k := range oldProps {
		if _, exists := newProps[k]; !exists {
			updates = append(updates, models.NewRemoveFilesetPropertyRequest(k))
		}
	}

	if len(updates) > 0 {
		result, err := r.client.UpdateFileset(state.Metalake.ValueString(), state.Catalog.ValueString(), state.Schema.ValueString(), state.Name.ValueString(), updates)
		if err != nil {
			resp.Diagnostics.AddError("Failed to update fileset", err.Error())
			return
		}
		setStateFromFileset(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), plan.Catalog.ValueString(), plan.Schema.ValueString(), &result.Fileset, &plan)
	} else {
		setStateFromFileset(ctx, &resp.Diagnostics, plan.Metalake.ValueString(), plan.Catalog.ValueString(), plan.Schema.ValueString(), nil, &plan)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FilesetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FilesetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DropFileset(state.Metalake.ValueString(), state.Catalog.ValueString(), state.Schema.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete fileset", err.Error())
		return
	}
}

func (r *FilesetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ".", 4)
	if len(parts) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected 'metalake.catalog.schema.fileset', got: %s", req.ID),
		)
		return
	}

	metalake := parts[0]
	catalog := parts[1]
	schemaName := parts[2]
	name := parts[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metalake"), metalake)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("catalog"), catalog)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schema"), schemaName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setStateFromFileset(ctx context.Context, diags *diag.Diagnostics, metalake, catalog, schemaName string, fileset *models.Fileset, model *FilesetResourceModel) {
	if fileset != nil {
		model.Name = types.StringValue(fileset.Name)
		model.Comment = types.StringValue(fileset.Comment)
		model.Type = types.StringValue(fileset.Type)
		model.StorageLocation = types.StringValue(fileset.StorageLocation)
		model.ID = types.StringValue(metalake + "." + catalog + "." + schemaName + "." + fileset.Name)

		props, d := types.MapValueFrom(ctx, types.StringType, fileset.Properties)
		diags.Append(d...)
		if !diags.HasError() {
			model.Properties = props
		}

		if fileset.Audit != nil {
			auditObj, d := auditToObjectValue(ctx, fileset.Audit)
			diags.Append(d...)
			if !diags.HasError() {
				model.Audit = auditObj
			}
		}
	} else {
		model.ID = types.StringValue(metalake + "." + catalog + "." + schemaName + "." + model.Name.ValueString())
	}

	model.Metalake = types.StringValue(metalake)
	model.Catalog = types.StringValue(catalog)
	model.Schema = types.StringValue(schemaName)
}

func auditToObjectValue(ctx context.Context, audit *models.Audit) (types.Object, diag.Diagnostics) {
	if audit == nil {
		return types.ObjectNull(AuditAttrTypes), nil
	}

	creator := types.StringValue(audit.Creator)
	lastModifier := types.StringValue(audit.LastModifier)

	var createTime, lastModifiedTime types.String
	if audit.CreateTime != nil {
		createTime = types.StringValue(audit.CreateTime.Format(time.RFC3339))
	} else {
		createTime = types.StringNull()
	}
	if audit.LastModifiedTime != nil {
		lastModifiedTime = types.StringValue(audit.LastModifiedTime.Format(time.RFC3339))
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
