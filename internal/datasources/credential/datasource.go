package credential

import (
	"context"
	"fmt"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ datasource.DataSource = &CredentialsDataSource{}
var _ datasource.DataSourceWithConfigure = &CredentialsDataSource{}

type CredentialsDataSource struct {
	client *client.Client
}

func New() datasource.DataSource {
	return &CredentialsDataSource{}
}

func (d *CredentialsDataSource) SetClient(c *client.Client) {
	d.client = c
}

type CredentialsDataSourceModel struct {
	Metalake     types.String `tfsdk:"metalake"`
	ResourceType types.String `tfsdk:"resource_type"`
	Resource     types.String `tfsdk:"resource"`
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
	ExpireTime   types.String `tfsdk:"expire_time"`
}

func (d *CredentialsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CredentialsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "gravitino_credentials"
}

func (d *CredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required:    true,
				Description: "The metalake name.",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The metadata object type (e.g. CATALOG, SCHEMA, TABLE, COLUMN, FILESET, TOPIC, MODEL, ROLE).",
				Validators: []validator.String{
					stringvalidator.OneOf(models.StatisticsObjectTypes...),
				},
			},
			"resource": schema.StringAttribute{
				Required:    true,
				Description: "The resource name.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The credential type.",
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The credential value.",
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "The credential expiration time.",
			},
		},
	}
}

func (d *CredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CredentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetCredentials(
		config.Metalake.ValueString(),
		config.ResourceType.ValueString(),
		config.Resource.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get credentials", err.Error())
		return
	}

	setCredentialState(ctx, &result.Credential, &config)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func setCredentialState(_ context.Context, cred *models.Credential, model *CredentialsDataSourceModel) {
	model.Type = types.StringValue(cred.Type)
	model.Value = types.StringValue(cred.Value)
	if cred.ExpireTime != "" {
		model.ExpireTime = types.StringValue(cred.ExpireTime)
	} else {
		model.ExpireTime = types.StringNull()
	}
}
