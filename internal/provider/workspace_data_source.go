package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &WorkspaceDataSource{}

func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

type WorkspaceDataSource struct {
	client *graphql.Client
}

type WorkspaceDataSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	UrlKey types.String `tfsdk:"url_key"`
}

func (d *WorkspaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (d *WorkspaceDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear workspace.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the workspace.",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				MarkdownDescription: "Name of the workspace.",
				Type:                types.StringType,
				Computed:            true,
			},
			"url_key": {
				MarkdownDescription: "URL key of the workspace.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *WorkspaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *WorkspaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspace(ctx, *d.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.Organization.Id)
	data.Name = types.StringValue(response.Organization.Name)
	data.UrlKey = types.StringValue(response.Organization.UrlKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
