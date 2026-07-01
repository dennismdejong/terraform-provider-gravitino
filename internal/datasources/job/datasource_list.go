package job

import (
	"context"
	"fmt"
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

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

var _ datasource.DataSource = &JobsDataSource{}
var _ datasource.DataSourceWithConfigure = &JobsDataSource{}

type JobsDataSource struct {
	client *client.Client
}

func NewListDataSource() datasource.DataSource {
	return &JobsDataSource{}
}

func (d *JobsDataSource) SetClient(c *client.Client) {
	d.client = c
}

type JobsDataSourceModel struct {
	Metalake types.String `tfsdk:"metalake"`
	Jobs     types.List   `tfsdk:"jobs"`
}

type jobItemModel struct {
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Schedule   types.String `tfsdk:"schedule"`
	Status     types.String `tfsdk:"status"`
	Parameters types.Map    `tfsdk:"parameters"`
	Audit      types.Object `tfsdk:"audit"`
}

var JobItemAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"template":   types.StringType,
	"schedule":   types.StringType,
	"status":     types.StringType,
	"parameters": types.MapType{ElemType: types.StringType},
	"audit":      types.ObjectType{AttrTypes: AuditAttrTypes},
}

func (d *JobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JobsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_jobs"
}

func (d *JobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"jobs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The job name.",
						},
						"template": schema.StringAttribute{
							Computed:    true,
							Description: "The job template name.",
						},
						"schedule": schema.StringAttribute{
							Computed:    true,
							Description: "The job schedule.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The current status of the job.",
						},
						"parameters": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The job parameters.",
						},
						"audit": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: AuditAttrTypes,
							Description:    "Audit information for the job.",
						},
					},
				},
			},
		},
	}
}

func (d *JobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JobsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListJobs(config.Metalake.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list jobs", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(result.Jobs))
	for _, job := range result.Jobs {
		j := job
		item := jobToItemModel(ctx, &j, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if item == nil {
			continue
		}
		obj, objDiags := types.ObjectValueFrom(ctx, JobItemAttrTypes, item)
		resp.Diagnostics.Append(objDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, obj)
	}

	jobsList, listDiags := types.ListValue(types.ObjectType{AttrTypes: JobItemAttrTypes}, items)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Jobs = jobsList
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func jobToItemModel(ctx context.Context, j *models.Job, diags *diag.Diagnostics) *jobItemModel {
	if j == nil {
		return nil
	}

	item := &jobItemModel{
		Name:     types.StringValue(j.Name),
		Template: types.StringValue(j.Template),
		Schedule: types.StringValue(j.Schedule),
		Status:   types.StringValue(j.Status),
	}

	props, d := types.MapValueFrom(ctx, types.StringType, strMapFromInterface(j.Parameters))
	if d.HasError() {
		return nil
	}
	item.Parameters = props

	auditObj, d := auditToObjectValueForDS(ctx, j.Audit)
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

func strMapFromInterface(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			result[k] = strVal
		}
	}
	return result
}
