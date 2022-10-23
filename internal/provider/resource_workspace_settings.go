package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = workspaceSettingsResourceType{}
var _ resource.Resource = workspaceSettingsResource{}
var _ resource.ResourceWithImportState = workspaceSettingsResource{}

type workspaceSettingsResourceType struct{}

func (t workspaceSettingsResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (t workspaceSettingsResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return workspaceSettingsResource{
		provider: provider,
	}, diags
}

type workspaceSettingsResourceData struct {
	Id                              types.String `tfsdk:"id"`
	EnableRoadmap                   types.Bool   `tfsdk:"enable_roadmap"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
}

type workspaceSettingsResource struct {
	provider linearProvider
}

func (r workspaceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data workspaceSettingsResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.Value,
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.Value,
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.Value,
	}

	response, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace settings, got error: %s", err))
		return
	}

	organization := response.OrganizationUpdate.Organization

	data.Id = types.String{Value: organization.Id}
	data.EnableRoadmap = types.Bool{Value: organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: organization.GitPublicLinkbackMessagesEnabled}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data workspaceSettingsResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspaceSettings(context.Background(), r.provider.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace settings, got error: %s", err))
		return
	}

	organization := response.Organization

	data.Id = types.String{Value: organization.Id}
	data.EnableRoadmap = types.Bool{Value: organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: organization.GitPublicLinkbackMessagesEnabled}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data workspaceSettingsResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.Value,
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.Value,
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.Value,
	}

	response, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data workspaceSettingsResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   false,
		GitLinkbackMessagesEnabled:       false,
		GitPublicLinkbackMessagesEnabled: false,
	}

	_, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted workspace settings")
}

func (r workspaceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := getWorkspaceSettings(context.Background(), r.provider.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace settings, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.Organization.Id)...)
}
