package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = workspaceDataSourceType{}
var _ datasource.DataSource = workspaceDataSource{}

type workspaceDataSourceType struct{}

func (t workspaceDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (t workspaceDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return workspaceDataSource{
		provider: provider,
	}, diags
}

type workspaceDataSourceData struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	UrlKey types.String `tfsdk:"url_key"`
}

type workspaceDataSource struct {
	provider linearProvider
}

func (d workspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspace(context.Background(), d.provider.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Organization.Id}
	data.Name = types.String{Value: response.Organization.Name}
	data.UrlKey = types.String{Value: response.Organization.UrlKey}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
