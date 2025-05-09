package provider

import (
	"context"
	"net/http"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Khan/genqlient/graphql"
)

var (
	envVarName          = "LINEAR_TOKEN"
	errMissingAuthToken = "Required token could not be found. Please set the token using an input variable in the provider configuration block or by using the `" + envVarName + "` environment variable."
)

func colorRegex() *regexp.Regexp {
	return regexp.MustCompile("^#[0-9a-fA-F]{6}$")
}

func uuidRegex() *regexp.Regexp {
	return regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
}

var _ provider.Provider = &LinearProvider{}

type LinearProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type LinearProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *LinearProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "linear"
	resp.Version = p.version
}

func (p *LinearProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "The token used to authenticate with Linear.",
				Optional:            true,
			},
		},
	}
}

func (p *LinearProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LinearProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	token := ""

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	// If a token wasn't set in the provider configuration block, try and fetch it
	// from the environment variable.
	if token == "" {
		token = os.Getenv(envVarName)
	}

	// If we still don't have a token at this point, we return an error.
	if token == "" {
		resp.Diagnostics.AddError("Missing API token", errMissingAuthToken)
		return
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			token:   token,
			wrapped: http.DefaultTransport,
		},
	}

	client := graphql.NewClient("https://api.linear.app/graphql", &httpClient)

	resp.DataSourceData = &client
	resp.ResourceData = &client
}

func (p *LinearProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTeamResource,
		NewTeamLabelResource,
		NewTeamWorkflowResource,
		NewTemplateResource,
		NewWorkflowStateResource,
		NewWorkspaceLabelResource,
		NewWorkspaceSettingsResource,
	}
}

func (p *LinearProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LinearProvider{
			version: version,
		}
	}
}
