package provider

import (
	"context"
	"os"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	dsauthentication "github.com/gravitino/terraform-provider-gravitino/internal/datasources/authentication"
	dscatalog "github.com/gravitino/terraform-provider-gravitino/internal/datasources/catalog"
	dscredential "github.com/gravitino/terraform-provider-gravitino/internal/datasources/credential"
	dsfileset "github.com/gravitino/terraform-provider-gravitino/internal/datasources/fileset"
	dsfunction "github.com/gravitino/terraform-provider-gravitino/internal/datasources/function"
	dsgroup "github.com/gravitino/terraform-provider-gravitino/internal/datasources/group"
	dshealth "github.com/gravitino/terraform-provider-gravitino/internal/datasources/health"
	dsjob "github.com/gravitino/terraform-provider-gravitino/internal/datasources/job"
	dsmetalake "github.com/gravitino/terraform-provider-gravitino/internal/datasources/metalake"
	dsmodel "github.com/gravitino/terraform-provider-gravitino/internal/datasources/model"
	dsowner "github.com/gravitino/terraform-provider-gravitino/internal/datasources/owner"
	dspartition "github.com/gravitino/terraform-provider-gravitino/internal/datasources/partition"
	dspolicy "github.com/gravitino/terraform-provider-gravitino/internal/datasources/policy"
	dsrole "github.com/gravitino/terraform-provider-gravitino/internal/datasources/role"
	dsschema "github.com/gravitino/terraform-provider-gravitino/internal/datasources/schema"
	dsstatistics "github.com/gravitino/terraform-provider-gravitino/internal/datasources/statistics"
	dstable "github.com/gravitino/terraform-provider-gravitino/internal/datasources/table"
	dstag "github.com/gravitino/terraform-provider-gravitino/internal/datasources/tag"
	dstopic "github.com/gravitino/terraform-provider-gravitino/internal/datasources/topic"
	dsuser "github.com/gravitino/terraform-provider-gravitino/internal/datasources/user"
	dsview "github.com/gravitino/terraform-provider-gravitino/internal/datasources/view"
	rscatalog "github.com/gravitino/terraform-provider-gravitino/internal/resources/catalog"
	rsfileset "github.com/gravitino/terraform-provider-gravitino/internal/resources/fileset"
	rsfunction "github.com/gravitino/terraform-provider-gravitino/internal/resources/function"
	rsgroup "github.com/gravitino/terraform-provider-gravitino/internal/resources/group"
	rsjob "github.com/gravitino/terraform-provider-gravitino/internal/resources/job"
	rsmetalake "github.com/gravitino/terraform-provider-gravitino/internal/resources/metalake"
	rsmodel "github.com/gravitino/terraform-provider-gravitino/internal/resources/model"
	rsowner "github.com/gravitino/terraform-provider-gravitino/internal/resources/owner"
	rspartition "github.com/gravitino/terraform-provider-gravitino/internal/resources/partition"
	rspolicy "github.com/gravitino/terraform-provider-gravitino/internal/resources/policy"
	rsrole "github.com/gravitino/terraform-provider-gravitino/internal/resources/role"
	rsschema "github.com/gravitino/terraform-provider-gravitino/internal/resources/schema"
	rstag "github.com/gravitino/terraform-provider-gravitino/internal/resources/tag"
	rstable "github.com/gravitino/terraform-provider-gravitino/internal/resources/table"
	rstopic "github.com/gravitino/terraform-provider-gravitino/internal/resources/topic"
	rsuser "github.com/gravitino/terraform-provider-gravitino/internal/resources/user"
	rsview "github.com/gravitino/terraform-provider-gravitino/internal/resources/view"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ provider.Provider = (*GravitinoProvider)(nil)

type GravitinoProvider struct {
	version string
}

type GravitinoProviderModel struct {
	URI        types.String `tfsdk:"uri"`
	Auth       types.String `tfsdk:"auth"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	OAuthToken types.String `tfsdk:"oauth_token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GravitinoProvider{
			version: version,
		}
	}
}

func (p *GravitinoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gravitino"
	resp.Version = p.version
}

func (p *GravitinoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Optional:    true,
				Description: "The URI of the Gravitino server. Can also be set via GRAVITINO_URI environment variable.",
			},
		"auth": schema.StringAttribute{
			Optional:    true,
			Description: "Authentication method: 'basic' or 'oauth'. Can also be set via GRAVITINO_AUTH environment variable.",
			Validators: []validator.String{
				stringvalidator.OneOf("basic", "oauth"),
			},
		},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username for basic authentication. Can also be set via GRAVITINO_USERNAME environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for basic authentication. Can also be set via GRAVITINO_PASSWORD environment variable.",
			},
			"oauth_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "OAuth2 bearer token. Can also be set via GRAVITINO_OAUTH_TOKEN environment variable.",
			},
		},
	}
}

func (p *GravitinoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config GravitinoProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uri := os.Getenv("GRAVITINO_URI")
	if !config.URI.IsNull() {
		uri = config.URI.ValueString()
	}
	if uri == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("uri"),
			"Missing Gravitino URI",
			"Set the uri attribute in the provider block or the GRAVITINO_URI environment variable.",
		)
		return
	}

	auth := os.Getenv("GRAVITINO_AUTH")
	if !config.Auth.IsNull() {
		auth = config.Auth.ValueString()
	}

	username := os.Getenv("GRAVITINO_USERNAME")
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	password := os.Getenv("GRAVITINO_PASSWORD")
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	oauthToken := os.Getenv("GRAVITINO_OAUTH_TOKEN")
	if !config.OAuthToken.IsNull() {
		oauthToken = config.OAuthToken.ValueString()
	}

	c, err := client.New(uri, auth, username, password, oauthToken)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create client", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *GravitinoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dsmetalake.NewMetalakesDataSource,
		dsmetalake.NewMetalakeDataSource,
		dscatalog.NewListDataSource,
		dscatalog.NewGetDataSource,
		dsjob.NewListDataSource,
		dsjob.NewGetDataSource,
		dsschema.NewSchemaDataSource,
		dsschema.NewSchemasDataSource,
		dsmodel.NewModelDataSource,
		dsmodel.NewModelsDataSource,
		dstag.NewListDataSource,
		dstag.NewGetDataSource,
		dsview.NewViewsDataSource,
		dsview.NewViewDataSource,
		dsfunction.NewFunctionsDataSource,
		dsfunction.NewFunctionDataSource,
		dsfileset.NewListDataSource,
		dsfileset.NewGetDataSource,
		dstopic.NewTopicsDataSource,
		dstopic.NewTopicDataSource,
		dsgroup.NewListDataSource,
		dsgroup.NewGetDataSource,
		dspartition.NewPartitionsDataSource,
		dspartition.NewPartitionDataSource,
		dstable.NewTablesDataSource,
		dstable.NewTableDataSource,
		dshealth.NewHealthDataSource,
		dshealth.NewLivenessDataSource,
		dshealth.NewReadinessDataSource,
		dspolicy.NewListDataSource,
		dscredential.New,
		dsrole.New,
		dsrole.NewRoleDataSource,
		dsrole.NewRolesListDataSource,
		dsstatistics.New,
		dsstatistics.NewPartition,
		dsauthentication.New,
		dsowner.NewOwnerDataSource,
		dsuser.NewListDataSource,
		dsuser.NewGetDataSource,
	}
}

func (p *GravitinoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		rsmetalake.NewMetalakeResource,
		rscatalog.New,
		rsfileset.New,
		rsjob.New,
		rsschema.NewSchemaResource,
		rsmodel.New,
		rsowner.NewOwnerResource,
		rstable.NewTableResource,
		rstag.New,
		rstopic.NewTopicResource,
		rsview.NewViewResource,
		rsfunction.New,
		rsgroup.New,
		rspartition.NewPartitionResource,
		rspolicy.New,
		rsrole.NewRoleResource,
		rsuser.New,
	}
}
