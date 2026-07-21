package job_template

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

var _ datasource.DataSource = &JobTemplatesDataSource{}
var _ datasource.DataSourceWithConfigure = &JobTemplatesDataSource{}

type JobTemplatesDataSource struct {
	client *client.Client
}

func NewJobTemplatesDataSource() datasource.DataSource {
	return &JobTemplatesDataSource{}
}

func (d *JobTemplatesDataSource) SetClient(c *client.Client) {
	d.client = c
}

type JobTemplatesDataSourceModel struct {
	Metalake  types.String `tfsdk:"metalake"`
	Templates types.List   `tfsdk:"templates"`
}

type jobTemplateItemModel struct {
	Name       types.String `tfsdk:"name"`
	Template   types.String `tfsdk:"template"`
	Parameters types.Map    `tfsdk:"parameters"`
	Comment    types.String `tfsdk:"comment"`
	Properties types.Map    `tfsdk:"properties"`
	Audit      types.Object `tfsdk:"audit"`
}

var JobTemplateItemAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"template":   types.StringType,
	"parameters": types.MapType{ElemType: types.StringType},
	"comment":    types.StringType,
	"properties": types.MapType{ElemType: types.StringType},
	"audit":      types.ObjectType{AttrTypes: AuditAttrTypes},
}

var AuditAttrTypes = map[string]attr.Type{
	"creator":            types.StringType,
	"create_time":        types.StringType,
	"last_modifier":      types.StringType,
	"last_modified_time": types.StringType,
}

func (d *JobTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JobTemplatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_job_templates"
}

func (d *JobTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The job template name.",
						},
						"template": schema.StringAttribute{
							Computed:    true,
							Description: "The template definition.",
						},
						"parameters": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The job template parameters.",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "The job template comment.",
						},
						"properties": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The job template properties.",
						},
						"audit": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: AuditAttrTypes,
							Description:    "Audit information for the job template.",
						},
					},
				},
			},
		},
	}
}

func (d *JobTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config JobTemplatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListJobTemplates(config.Metalake.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list job templates", err.Error())
		return
	}

	items := make([]attr.Value, 0, len(result.JobTemplates))
	for _, jt := range result.JobTemplates {
		t := jt
		item := jobTemplateToItemModel(ctx, &t, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if item == nil {
			continue
		}
		obj, objDiags := types.ObjectValueFrom(ctx, JobTemplateItemAttrTypes, item)
		resp.Diagnostics.Append(objDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, obj)
	}

	templatesList, listDiags := types.ListValue(types.ObjectType{AttrTypes: JobTemplateItemAttrTypes}, items)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Templates = templatesList
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func jobTemplateToItemModel(ctx context.Context, jt *models.JobTemplate, diags *diag.Diagnostics) *jobTemplateItemModel {
	if jt == nil {
		return nil
	}

	item := &jobTemplateItemModel{
		Name:     types.StringValue(jt.Name),
		Template: types.StringValue(jt.Template),
		Comment:  types.StringValue(jt.Comment),
	}

	params, d := types.MapValueFrom(ctx, types.StringType, jt.Parameters)
	if d.HasError() {
		return nil
	}
	item.Parameters = params

	props, d := types.MapValueFrom(ctx, types.StringType, jt.Properties)
	if d.HasError() {
		return nil
	}
	item.Properties = props

	auditObj, d := auditToObjectValueForDS(ctx, jt.Audit)
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

	attrs := map[string]attr.Value{
		"creator":       types.StringValue(audit.Creator),
		"last_modifier": types.StringValue(audit.LastModifier),
	}

	if audit.CreateTime != nil {
		attrs["create_time"] = types.StringValue(audit.CreateTime.Format(time.RFC3339))
	} else {
		attrs["create_time"] = types.StringNull()
	}
	if audit.LastModifiedTime != nil {
		attrs["last_modified_time"] = types.StringValue(audit.LastModifiedTime.Format(time.RFC3339))
	} else {
		attrs["last_modified_time"] = types.StringNull()
	}

	return types.ObjectValue(AuditAttrTypes, attrs)
}
