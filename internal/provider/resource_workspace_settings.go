package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
)

var _ resource.Resource = &WorkspaceSettingsResource{}
var _ resource.ResourceWithImportState = &WorkspaceSettingsResource{}

func NewWorkspaceSettingsResource() resource.Resource {
	return &WorkspaceSettingsResource{}
}

type WorkspaceSettingsResource struct {
	client *graphql.Client
}

type WorkspaceSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	EnableRoadmap                   types.Bool   `tfsdk:"enable_roadmap"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
}

func (r *WorkspaceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_settings"
}

func (r *WorkspaceSettingsResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear workspace settings.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the workspace.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"enable_roadmap": {
				MarkdownDescription: "Enable roadmap for the workspace. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(false),
				},
			},
			"enable_git_linkback_messages": {
				MarkdownDescription: "Enable git linkbacks for private repositories. **Default** `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(true),
				},
			},
			"enable_git_linkback_messages_public": {
				MarkdownDescription: "Enable git linkbacks for public repositories. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(false),
				},
			},
		},
	}, nil
}

func (r *WorkspaceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *WorkspaceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.Value,
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.Value,
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.Value,
	}

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace settings, got error: %s", err))
		return
	}

	organization := response.OrganizationUpdate.Organization

	data.Id = types.String{Value: organization.Id}
	data.EnableRoadmap = types.Bool{Value: organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: organization.GitPublicLinkbackMessagesEnabled}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspaceSettings(ctx, *r.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace settings, got error: %s", err))
		return
	}

	organization := response.Organization

	data.Id = types.String{Value: organization.Id}
	data.EnableRoadmap = types.Bool{Value: organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: organization.GitPublicLinkbackMessagesEnabled}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.Value,
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.Value,
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.Value,
	}

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated workspace settings")

	organization := response.OrganizationUpdate.Organization

	data.Id = types.String{Value: organization.Id}
	data.EnableRoadmap = types.Bool{Value: organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: organization.GitPublicLinkbackMessagesEnabled}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   false,
		GitLinkbackMessagesEnabled:       false,
		GitPublicLinkbackMessagesEnabled: false,
	}

	_, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted workspace settings")
}

func (r *WorkspaceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := getWorkspaceSettings(ctx, *r.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace settings, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.Organization.Id)...)
}
